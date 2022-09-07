package transaction

import "crypto/ecdsa"

type blockchain interface {
	FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int)
	SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey)
}
