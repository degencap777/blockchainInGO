package cli

import (
	"fmt"

	"github.com/noodleslove/blockchain-go/pkg/utils"
	"github.com/noodleslove/blockchain-go/pkg/wallet"
)

func (cli *CLI) listAddresses() {
	wallets, err := wallet.NewWallets()
	utils.Check(err)
	addresses := wallets.GetAddresses()

	for _, address := range addresses {
		fmt.Println(address)
	}
}
