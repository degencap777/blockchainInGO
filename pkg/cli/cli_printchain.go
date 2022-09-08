package cli

import (
	"fmt"
	"strconv"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
)

func (cli *CLI) printChain() {
	bc := blockchain.NewBlockchain()
	defer bc.CloseDB()
	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
