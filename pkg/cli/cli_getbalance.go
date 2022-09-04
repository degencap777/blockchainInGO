package cli

import (
	"fmt"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
)

func (cli *CLI) getBalance(address string) {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
