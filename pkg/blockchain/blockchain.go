package main

type Blockchain struct {
	blocks []*Block
}

// AddBlock saves the provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlockHash := bc.blocks[len(bc.blocks)-1].Hash
	newBlock := NewBlock(data, prevBlockHash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewBlockchain returns a blockchain with a genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{
		blocks: []*Block{NewGenesisBlock()},
	}
}
