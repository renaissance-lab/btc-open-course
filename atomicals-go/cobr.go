package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"

	cbor "github.com/fxamacker/cbor/v2"
)

type AtomPayLoad struct {
	Args AtomArg `cbor:"args"`
}

type AtomArg struct {
	Bitworkc   string `cbor:"bitworkc"`
	Bitworkr   string `cbor:"bitworkr,omitempty"`
	MintTicker string `cbor:"mint_ticker"`
	Nonce      int64  `cbor:"nonce"`
	Time       int64  `cbor:"time"`
}

func PayloadToCbor(payLoad AtomPayLoad) []byte {

	b, err := cbor.Marshal(payLoad)
	if err != nil {
		fmt.Println("encode error:", err)
	}

	return b
}

func CborToPayLoad(payLoadBytes string) {

	expectData := payLoadBytes

	expectDataBytes, _ := hex.DecodeString(expectData)

	dec := cbor.NewDecoder(bytes.NewReader(expectDataBytes))
	for {
		var animal AtomPayLoad
		if err := dec.Decode(&animal); err != nil {
			if err != io.EOF {
				fmt.Println("decode error:", err)
			}
			break
		}
		fmt.Printf("%+v\n", animal)
	}

}
