package main

import "math/big"

const targetBits = 18

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds and returns a new proofofwork consensus
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	return &ProofOfWork{
		block:  b,
		target: target,
	}
}
