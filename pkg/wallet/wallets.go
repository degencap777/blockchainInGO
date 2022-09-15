package wallet

import (
	"bytes"
	"crypto/x509"
	"encoding/gob"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/noodleslove/blockchain-go/internal"
	"github.com/noodleslove/blockchain-go/pkg/utils"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

type encodedWallets map[string][2][]byte

type encodedWallet [2][]byte

func NewWallets(nodeID string) (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadFromFile(nodeID)

	return &wallets, err
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := string(wallet.GetAddress())

	ws.Wallets[address] = wallet

	return address
}

// GetWallet returns a Wallet by its address
func (ws *Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

// GetAddresses returns an array of addresses stored in the wallet file
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

// LoadFromFile loads wallets from the file
func (ws *Wallets) LoadFromFile(nodeID string) error {
	walletFile := fmt.Sprintf(internal.WalletFile, nodeID)
	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		return err
	}

	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	var wallets map[string][2][]byte
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		return err
	}

	ws.Decode(wallets)

	return nil
}

// SaveToFile saves wallets to a file
func (ws *Wallets) SaveToFile(nodeID string) {
	var content bytes.Buffer
	walletFile := fmt.Sprintf(internal.WalletFile, nodeID)

	wallets := ws.Encode()
	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(wallets)
	utils.Check(err)

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	utils.Check(err)
}

func (ws *Wallets) Encode() encodedWallets {
	encodedWallets := make(encodedWallets)

	for addr, w := range ws.Wallets {
		x509Encoded, _ := x509.MarshalECPrivateKey(&w.PrivateKey)
		pemEncoded := pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509Encoded,
		})

		encodedWallets[addr] = encodedWallet{pemEncoded, w.PublicKey}
	}

	return encodedWallets
}

func (ws *Wallets) Decode(encodedWallets encodedWallets) {
	for addr, w := range encodedWallets {
		encPriv, encPub := w[0], w[1]

		block, _ := pem.Decode([]byte(encPriv))
		x509Encoded := block.Bytes
		privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

		ws.Wallets[addr] = &Wallet{
			PrivateKey: *privateKey,
			PublicKey:  encPub,
		}
	}
}
