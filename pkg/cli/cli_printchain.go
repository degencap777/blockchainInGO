package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
)

func (cli *CLI) printChain(nodeID string) {
	bc := blockchain.NewBlockchain(nodeID)
	defer bc.CloseDB()
	bci := bc.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("============ Block %x ============\n", block.Hash)
		fmt.Printf("Height: %d\n", block.Height)
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Created at : %s\n", time.Unix(block.Timestamp, 0))
		fmt.Printf("Merkle Root: %x\n", block.HashTransactions())
		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("PoW: %s\n\n", strconv.FormatBool(pow.Validate()))
		for _, tx := range block.Transactions {
			fmt.Println(tx)
		}
		fmt.Printf("\n\n")

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}
