package cli

import (
	"fmt"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
)

func (cli *CLI) createBlockchain(address string) {
	bc := blockchain.CreateBlockchain(address)
	bc.CloseDB()
	fmt.Println("Done!")
}
