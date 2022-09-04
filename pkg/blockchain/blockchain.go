package blockchain

import (
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
		cbtx := transaction.NewCoinbaseTx(address, genesisData)
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

func (bc *Blockchain) CloseDB() {
	bc.db.Close()
}
