package cli

import (
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"github.com/noodleslove/blockchain-go/internal"
	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/wallet"
)

func (cli *CLI) getBalance(address, nodeID string) {
	if !wallet.ValidateAddress(address) {
		log.Panic("ERROR: Address is not valid")
	}

	bc := blockchain.NewBlockchain(nodeID)
	utxoSet := blockchain.UTXOSet{Blockchain: bc}
	defer bc.CloseDB()

	balance := 0
	pubKeyHash := base58.Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-internal.AddressChecksumLen]
	UTXOs := utxoSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}
