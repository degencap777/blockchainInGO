package network

import (
	"bytes"
	"encoding/gob"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type getblocks struct {
	AddrFrom string
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func handleGetBlocks(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	// First, we need to decode the request and extract the payload
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.Check(err)

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}
