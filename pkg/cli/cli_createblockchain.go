package cli

import (
	"fmt"
	"log"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/wallet"
)

func (cli *CLI) createBlockchain(address, nodeID string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := blockchain.CreateBlockchain(address, nodeID)
	defer bc.CloseDB()

	utxoSet := blockchain.UTXOSet{Blockchain: bc}
	utxoSet.Reindex()

	fmt.Println("Done!")
}
