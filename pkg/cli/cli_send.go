package cli

import (
	"fmt"
	"log"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/network"
	"github.com/noodleslove/blockchain-go/pkg/wallet"
)

func (cli *CLI) send(from, to string, amount int, nodeID string, mineNow bool) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := blockchain.NewBlockchain(nodeID)
	utxoSet := blockchain.UTXOSet{Blockchain: bc}
	defer bc.CloseDB()

	tx := blockchain.NewUTXOTransaction(from, to, amount, utxoSet)

	if mineNow {
		cbTx := blockchain.NewCoinbaseTX(from, "")
		txs := []*blockchain.Transaction{cbTx, tx}

		newBlock := bc.MineBlock(txs)
		utxoSet.Update(newBlock)
	} else {
		network.SendTx(network.KnownNodes[0], tx)
	}
	fmt.Println("Success!")
}
