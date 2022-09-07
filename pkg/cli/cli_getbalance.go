package cli

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/noodleslove/blockchain-go/pkg/blockchain"
)

func (cli *CLI) getBalance(address string) {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	balance := 0
	pubKeyHash := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4] // TODO: fix checksumLen
	UTXOs := bc.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
