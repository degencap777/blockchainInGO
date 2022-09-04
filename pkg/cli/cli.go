package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type CLI struct {
	bc *blockchain.Blockchain
}

func NewCLI(bc *blockchain.Blockchain) *CLI {
	return &CLI{
		bc: bc,
	}
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) Run() {
	cli.validateArgs()

	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	switch os.Args[1] {
	case "addblock":
		err := sendCmd.Parse(os.Args[2:])
		utils.Check(err)
	case "printchain":
		err := printChainCmd.Parse(os.Args[2:])
		utils.Check(err)
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
