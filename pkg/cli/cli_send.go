package cli

import (
	"fmt"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/transaction"
)

func (cli *CLI) send(from, to string, amount int) {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	tx := transaction.NewUTXOTransaction(from, to, amount, bc)
	cbTx := transaction.NewCoinbaseTX(from, "")
	txs := []*transaction.Transaction{cbTx, tx}

	bc.MineBlock(txs)
	fmt.Println("Success!")
}
