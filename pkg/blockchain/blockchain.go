package blockchain

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
	"github.com/noodleslove/blockchain-go/pkg/transaction"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"
const genesisData = "Bitcoin was created in 2009 by a person or group of people using the pseudonym Satoshi Nakamoto"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock saves the provided data as a block in the blockchain
func (bc *Blockchain) MineBlock(transactions []*transaction.Transaction) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		lastHash = b.Get([]byte("l"))

		return nil
	})
	utils.Check(err)

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		utils.Check(err)
		err = b.Put([]byte("l"), newBlock.Hash)
		utils.Check(err)
		bc.tip = newBlock.Hash

		return nil
	})
	utils.Check(err)
}

// Helper function check if blockchain db exists
func dbExists() bool {
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
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
	db, err := bolt.Open(dbFile, 0600, nil)
	utils.Check(err)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
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
	db, err := bolt.Open(dbFile, 0600, nil)
	utils.Check(err)

	err = db.Update(func(tx *bolt.Tx) error {
		cbtx := transaction.NewCoinbaseTX(address, genesisData)
		genesis := NewGenesisBlock(cbtx)
		b, err := tx.CreateBucket([]byte(blockBucket))
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
func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []transaction.Transaction {
	var unspentTXs []transaction.Transaction
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

// FindUTXO finds and returns all unspent transaction outputs
func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []transaction.TXOutput {
	var UTXOs []transaction.TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

func (bc *Blockchain) FindSpendableOutputs(
	pubKeyHash []byte,
	amount int,
) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)
	accumulated := 0

Outputs:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Outputs
				}
			}
		}
	}

	return accumulated, unspentOutputs
}
