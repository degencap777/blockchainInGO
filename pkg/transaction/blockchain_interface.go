package transaction

import "crypto/ecdsa"

type blockchain interface {
	SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey)
}
