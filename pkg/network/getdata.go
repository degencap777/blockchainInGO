package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type getdata struct {
	AddrFrom string
	Type     string
	ID       []byte
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{
		AddrFrom: nodeAddress,
		Type:     kind,
		ID:       id,
	})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func handleGetData(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	// First, we need to decode the request and extract the payload
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.Check(err)

	if payload.Type == "block" {
		block, err := bc.GetBlock([]byte(payload.ID))
		utils.Check(err)

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		SendTx(payload.AddrFrom, &tx)
	}
}
