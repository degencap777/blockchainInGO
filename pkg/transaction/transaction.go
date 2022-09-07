package transaction

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/noodleslove/blockchain-go/pkg/utils"
)

const subsidy = 10

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

// NewCoinbaseTX creates a new coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txin := TXInput{
		Txid:      []byte{},
		Vout:      -1,
		ScriptSig: data,
	}
	txout := TXOutput{
		Value:        subsidy,
		ScriptPubKey: to,
	}
	tx := Transaction{
		ID:   nil,
		Vin:  []TXInput{txin},
		Vout: []TXOutput{txout},
	}
	tx.SetID()

	return &tx
}

// SetID sets ID of a transaction
func (tx *Transaction) SetID() {
	var encoded bytes.Buffer
	var hash [32]byte

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	utils.Check(err)
	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = hash[:]
}

// IsCoinbase determines if a transaction is coinbase
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

func NewUTXOTransaction(from, to string, amount int, bc blockchain) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panic("ERROR: Not enough funds")
	}

	// Build a list of inputs
	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		utils.Check(err)

		for _, out := range outs {
			input := TXInput{txID, out, from}
			inputs = append(inputs, input)
		}
	}

	// Build a list of outputs
	outputs = append(outputs, TXOutput{amount, to})
	if acc > amount {
		outputs = append(outputs, TXOutput{acc - amount, from}) // a change
	}

	tx := Transaction{
		ID:   nil,
		Vin:  inputs,
		Vout: outputs,
	}
	tx.SetID()

	return &tx
}

func (tx *Transaction) Sign(privKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	// Coinbase transactions are not signed because there are no real inputs in them.
	if tx.IsCoinbase() {
		return
	}

	// A trimmed copy will be signed, not a full transaction.
	txCopy := tx.TrimmedCopy()

	// Next, we iterate over each input in the copy.
	for inID, vin := range txCopy.Vin {
		// In each input, Signature is set to nil (just a double-check) and
		// PubKey is set to the PubKeyHash of the referenced output.
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash

		// The Hash method serializes the transaction and hashes it with the
		// SHA-256 algorithm. The resulted hash is the data we’re going to sign.
		txCopy.ID = txCopy.Hash()

		// After getting the hash we should reset the PubKey field, so it doesn’t
		// affect further iterations.
		txCopy.Vin[inID].PubKey = nil

		// We sign txCopy.ID with privKey. An ECDSA signature is a pair of numbers,
		// which we concatenate and store in the input’s Signature field.
		r, s, err := ecdsa.Sign(rand.Reader, &privKey, txCopy.ID)
		utils.Check(err)
		signature := append(r.Bytes(), s.Bytes()...)

		tx.Vin[inID].Signature = signature
	}
}

// TrimmedCopy creates a trimmed copy of Transaction to be used in signing
func (tx *Transaction) TrimmedCopy() Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	for _, vin := range tx.Vin {
		inputs = append(inputs, TXInput{vin.Txid, vin.Vout, nil, nil})
	}

	for _, vout := range tx.Vout {
		outputs = append(outputs, TXOutput{vout.Value, vout.PubKeyHash})
	}

	txCopy := Transaction{tx.ID, inputs, outputs}

	return txCopy
}

// Hash returns the hash of the Transaction
func (tx *Transaction) Hash() []byte {
	var hash [32]byte

	txCopy := *tx
	txCopy.ID = []byte{}

	hash = sha256.Sum256(txCopy.Serialize())

	return hash[:]
}

// Serialize returns a serialized Transaction
func (tx *Transaction) Serialize() []byte {
	var encoded bytes.Buffer

	enc := gob.NewEncoder(&encoded)
	err := enc.Encode(tx)
	utils.Check(err)

	return encoded.Bytes()
}

func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	// First, we need the same transaction copy.
	txCopy := tx.TrimmedCopy()
	// Next, we’ll need the same curve that is used to generate key pairs.
	curve := elliptic.P256()

	// Next, we check signature in each input.
	for inID, vin := range tx.Vin {
		prevTx := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].Signature = nil
		txCopy.Vin[inID].PubKey = prevTx.Vout[vin.Vout].PubKeyHash
		txCopy.ID = txCopy.Hash()
		txCopy.Vin[inID].PubKey = nil

		// Here we unpack values stored in TXInput.Signature and TXInput.PubKey,
		// since a signature is a pair of numbers and a public key is a pair of
		// coordinates. We concatenated them earlier for storing, and now we need
		// to unpack them to use in crypto/ecdsa functions.
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		// Here it is: we create an ecdsa.PublicKey using the public key extracted
		// from the input and execute ecdsa.Verify passing the signature extracted
		// from the input. If all inputs are verified, return true; if at least
		// one input fails verification, return false.
		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if !ecdsa.Verify(&rawPubKey, txCopy.ID, &r, &s) {
			return false
		}
	}

	return true
}
