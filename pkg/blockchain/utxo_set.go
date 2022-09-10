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
