package blockchain

import (
	"encoding/hex"

	"github.com/boltdb/bolt"
	"github.com/noodleslove/blockchain-go/internal"
	"github.com/noodleslove/blockchain-go/pkg/transaction"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type UTXOSet struct {
	Blockchain *Blockchain
}

// Reindex rebuilds the UTXO set
func (u UTXOSet) Reindex() {
	db := u.Blockchain.db
	bucketName := []byte(internal.UtxoBucket)

	err := db.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucketName)
		utils.Check(err)
		_, err = tx.CreateBucket(bucketName)

		return err
	})
	utils.Check(err)

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			utils.Check(err)
			err = b.Put(key, outs.Serialize())
			utils.Check(err)
		}

		return nil
	})
	utils.Check(err)
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (u *UTXOSet) FindSpendableOutputs(
	pubKeyHash []byte,
	amount int,
) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumlated := 0
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(internal.UtxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			txID := hex.EncodeToString(k)
			outs := transaction.DeserializeOutputs(v)

			for outIdx, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) && accumlated < amount {
					accumlated += out.Value
					unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
				}
			}
		}

		return nil
	})
	utils.Check(err)

	return accumlated, unspentOutputs

}

// FindUTXO finds UTXO for a public key hash
func (u *UTXOSet) FindUTXO(pubKeyHash []byte) []transaction.TXOutput {
	var UTXOs []transaction.TXOutput
	db := u.Blockchain.db

	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(internal.UtxoBucket))
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			outs := transaction.DeserializeOutputs(v)

			for _, out := range outs.Outputs {
				if out.IsLockedWithKey(pubKeyHash) {
					UTXOs = append(UTXOs, out)
				}
			}
		}

		return nil
	})
	utils.Check(err)

	return UTXOs
}

// Update updates the UTXO set with transactions from the Block
// The Block is considered to be the tip of a blockchain
func (u *UTXOSet) Update(block *Block) {
	db := u.Blockchain.db

	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(internal.UtxoBucket))

		for _, tx := range block.Transactions {
			if !tx.IsCoinbase() {
				for _, vin := range tx.Vin {
					updateOuts := transaction.TXOutputs{}
					outsBytes := b.Get(vin.Txid)
					outs := transaction.DeserializeOutputs(outsBytes)

					// Put unspent outputs into updateOuts
					for outIdx, out := range outs.Outputs {
						if outIdx != vin.Vout {
							updateOuts.Outputs = append(updateOuts.Outputs, out)
						}
					}

					// Remove pair if all outputs are spent, otherwise save the
					// updated one
					if len(updateOuts.Outputs) == 0 {
						err := b.Delete(vin.Txid)
						utils.Check(err)
					} else {
						err := b.Put(vin.Txid, updateOuts.Serialize())
						utils.Check(err)
					}
				}
			}

			// Insert outputs of newly mined transactions
			newOutputs := transaction.TXOutputs{}
			for _, out := range tx.Vout {
				newOutputs.Outputs = append(newOutputs.Outputs, out)
			}

			err := b.Put(tx.ID, newOutputs.Serialize())
			utils.Check(err)
		}

		return nil
	})
	utils.Check(err)
}
