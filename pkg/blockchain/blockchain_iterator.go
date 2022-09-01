package main

import (
	"github.com/boltdb/bolt"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// Next returns next block starting from the tip
func (bci *BlockchainIterator) Next() *Block {
	var block *Block

	err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blockBucket))
		encodedBlock := b.Get([]byte(bci.currentHash))
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	utils.Check(err)

	bci.currentHash = block.PrevBlockHash
	return block
}
