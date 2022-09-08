package cli

import (
	"fmt"

	"github.com/noodleslove/blockchain-go/pkg/wallet"
)

func (cli *CLI) createWallet() {
	wallets, _ := wallet.NewWallets()
	address := wallets.CreateWallet()
	wallets.SaveToFile()

	fmt.Printf("Your new address: %s\n", address)
}
