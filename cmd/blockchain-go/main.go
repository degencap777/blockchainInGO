package main

import (
	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/cli"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	cli := cli.NewCLI(bc)
	cli.Run()
}
