package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"

	"golang.org/x/crypto/ripemd160"
)

func Check(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// IntToHex converts num and returns a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	Check(err)
	return buff.Bytes()
}

// Helper function hashes public key
func HashPubKey(pubKey []byte) []byte {
	publicSHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	Check(err)
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160
}
