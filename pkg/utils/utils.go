package utils

import (
	"bytes"
	"encoding/binary"
	"log"
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
