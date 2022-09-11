package blockchain

import (
	"bytes"

	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := utils.HashPubKey(in.PubKey)

	return bytes.Equal(lockingHash, pubKeyHash)
}
