package blockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
	"github.com/noodleslove/blockchain-go/internal"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock saves the provided data as a block in the blockchain
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	var lastHash []byte
	var lastHeight int

	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(internal.BlockBucket))
		lastHash = b.Get([]byte("l"))

		blockData := b.Get(lastHash)
		block := DeserializeBlock(blockData)

		lastHeight = block.Height

		return nil
	})
	utils.Check(err)

	newBlock := NewBlock(transactions, lastHash, lastHeight+1)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(internal.BlockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		utils.Check(err)
		err = b.Put([]byte("l"), newBlock.Hash)
		utils.Check(err)
		bc.tip = newBlock.Hash

		return nil
	})
	utils.Check(err)

	return newBlock
}

// Helper function check if blockchain db exists
func dbExists() bool {
	if _, err := os.Stat(internal.DbFile); os.IsNotExist(err) {
		return false
	}

	return true
}

// NewBlockchain returns a blockchain from existing db
func NewBlockchain() *Blockchain {
	if !dbExists() {
		fmt.Println("No existing blockchain found. Create one first.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(internal.DbFile, 0600, nil)
	utils.Check(err)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(internal.BlockBucket))
		tip = b.Get([]byte("l"))

		return nil
	})
	utils.Check(err)

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}

// CreateBlockchain returns a blockchain with a genesis block
func CreateBlockchain(address string) *Blockchain {
	if dbExists() {
		fmt.Println("Blockchain already exists.")
		os.Exit(1)
	}

	var tip []byte
	db, err := bolt.Open(internal.DbFile, 0600, nil)
	utils.Check(err)

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := NewCoinbaseTX(address, internal.GenesisData)
		genesis := NewGenesisBlock(cbtx)
		b, err := tx.CreateBucket([]byte(internal.BlockBucket))
		utils.Check(err)
		err = b.Put([]byte(genesis.Hash), genesis.Serialize())
		utils.Check(err)
		err = b.Put([]byte("l"), genesis.Hash)
		utils.Check(err)
		tip = genesis.Hash

		return nil
	})
	utils.Check(err)

	return &Blockchain{
		tip: tip,
		db:  db,
	}
}

// Iterator generates an iterator
func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
}

// Helper function close blockchain db
func (bc *Blockchain) CloseDB() {
	bc.db.Close()
}

// FindUnspentTransactions returns a list of transactions containing unspent outputs
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				// Check if an output was already referenced in an input
				if spentTXOs[txID] != nil {
					for _, spentOut := range spentTXOs[txID] {
						if spentOut == outIdx {
							continue Outputs
						}
					}
				}

				// If an output was locked by the same pubkey hash, this is the
				// output we want
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			// After checking outputs we gather all inputs that could unlock
			// outputs locked with the provided address (this doesn't apply to
			// coinbase transactions, since they don't unlock outputs)
			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return unspentTXs
}

// FindUTXO finds all unspent transaction outputs and returns transactions with
// spent outputs removed
func (bc *Blockchain) FindUTXO() map[string]TXOutputs {
	UTXO := make(map[string]TXOutputs)
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for outIdx, out := range tx.Vout {
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outputs = append(outs.Outputs, out)
				UTXO[txID] = outs
			}

			if !tx.IsCoinbase() {
				for _, in := range tx.Vin {
					inTxID := hex.EncodeToString(in.Txid)
					spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return UTXO
}

func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Equal(tx.ID, ID) {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("transaction not found")
}

func (bc *Blockchain) SignTransaction(
	tx *Transaction,
	privKey ecdsa.PrivateKey,
) {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTx, err := bc.FindTransaction(vin.Txid)
		utils.Check(err)
		prevTXs[hex.EncodeToString(prevTx.ID)] = prevTx
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		utils.Check(err)
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}

// GetBestHeight returns the height of the latest block
func (bc *Blockchain) GetBestHeight() int {
	var lastBlock Block

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(internal.BlockBucket))
		lastHash := b.Get([]byte("l"))
		blockData := b.Get(lastHash)
		lastBlock = *DeserializeBlock(blockData)

		return nil
	})
	utils.Check(err)

	return lastBlock.Height
}
