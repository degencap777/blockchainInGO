package cli

import (
	"fmt"
	"log"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/transaction"
	"github.com/noodleslove/blockchain-go/pkg/wallet"
)

func (cli *CLI) send(from, to string, amount int) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	tx := transaction.NewUTXOTransaction(from, to, amount, bc)
	cbTx := transaction.NewCoinbaseTX(from, "")
	txs := []*transaction.Transaction{cbTx, tx}

	bc.MineBlock(txs)
	fmt.Println("Success!")
}
