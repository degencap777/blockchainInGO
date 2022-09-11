package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/noodleslove/blockchain-go/pkg/merkletree"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type Block struct {
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// NewBlock creates and returns a block
func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Hash:          []byte{},
		Nonce:         0,
	}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}

// NewBlock creates and returns a genesis block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize serializes a block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer

	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	utils.Check(err)

	return result.Bytes()
}

// DeserializeBlock deserializes a block
func DeserializeBlock(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&block)
	utils.Check(err)

	return &block
}

func (b *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range b.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	mTree := merkletree.NewMerkleTree(transactions)

	return mTree.RootNode.Data
}
