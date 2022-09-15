package network

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"fmt"

	"github.com/noodleslove/blockchain-go/pkg/blockchain"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type inv struct {
	AddrFrom string
	Type     string
	Items    [][]byte
}

func sendInv(address, kind string, items [][]byte) {
	payload := gobEncode(inv{
		AddrFrom: nodeAddress,
		Type:     kind,
		Items:    items,
	})
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func handleInv(request []byte, bc *blockchain.Blockchain) {
	var buff bytes.Buffer
	var payload inv

	// First, we need to decode the request and extract the payload
	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	utils.Check(err)

	fmt.Printf("Recevied inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		// Right after putting blocks into the transit state, we send getdata
		// command to the sender of the inv message and update blocksInTransit
		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		newInTransit := [][]byte{}
		for _, b := range blocksInTransit {
			if !bytes.Equal(b, blockHash) {
				newInTransit = append(newInTransit, b)
			}
		}
		blocksInTransit = newInTransit
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}
