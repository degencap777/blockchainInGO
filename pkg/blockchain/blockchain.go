package main

import (
	"github.com/boltdb/bolt"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

const dbFile = "blockchain.db"
const blockBucket = "blocks"

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

// AddBlock saves the provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlockHash := bc.blocks[len(bc.blocks)-1].Hash
	newBlock := NewBlock(data, prevBlockHash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewBlockchain returns a blockchain with a genesis block
func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	utils.Check(err)

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blockBucket))
			utils.Check(err)
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put("l", genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{
		tip: tip,
		db:  db,
	}
	return &bc
}
