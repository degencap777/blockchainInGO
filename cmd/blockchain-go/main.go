package main

import (
	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	cli_ "github.com/noodleslove/blockchain-go/pkg/cli"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()

	cli := cli_.NewCLI(bc)
	cli.Run()
}
