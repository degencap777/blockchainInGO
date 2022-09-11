package cli

import (
	"fmt"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
)

func (cli *CLI) reindexUTXO() {
	bc := blockchain.NewBlockchain()
	utxoSet := blockchain.UTXOSet{Blockchain: bc}
	utxoSet.Reindex()
	defer bc.CloseDB()

	count := utxoSet.CountTransactions()
	fmt.Printf("Done! There are %d transactions in the UTXO set.\n", count)
}
