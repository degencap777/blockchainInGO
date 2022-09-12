package network

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type block struct {
	AddrFrom string
	Block    []byte
}

func sendBlock(address string, b *blockchain.Block) {
	payload := gobEncode(block{
		AddrFrom: nodeAddress,
		Block:    b.Serialize(),
	})
	request := append(commandToBytes("block"), payload...)

	sendData(address, request)
}

func handleBlock(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload block

	// First, we need to decode the request and extract the payload
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.Check(err)

	blockData := payload.Block
	block := blockchain.DeserializeBlock(blockData)

	fmt.Println("Recevied a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := blockchain.UTXOSet{Blockchain: bc}
		UTXOSet.Reindex()
	}
}
