package blockchain

import (
	"encoding/hex"

	"github.com/boltdb/bolt"
	"github.com/noodleslove/blockchain-go/internal"
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
		_, err = tx.CreateBucket(bucketName)

		return err
	})
	utils.Check(err)

	UTXO := u.Blockchain.FindUTXO()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)

		for txID, outs := range UTXO {
			key, err := hex.DecodeString(txID)
			err = b.Put(key, outs.Serialize())

			utils.Check(err)
		}

		return nil
	})
}
