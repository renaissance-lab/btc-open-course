package main

import (
	"encoding/hex"
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
)

func P2SHP2WPHExamaple() {
	priv := []byte("btc-example-data")

	net := &chaincfg.TestNet3Params

	btcecPriKey, btcecPublicKey := btcec.PrivKeyFromBytes(priv)

	expectAddr := "2NBp8EeKz3FQjYDMxmoqdvHzzgLQeeZpLsM"

	wif, err := btcutil.NewWIF(btcecPriKey, net, true)
	if err != nil {
		fmt.Println("NewWIF err ", err)
		return
	}
	// privkey := wif.PrivKey
	pubkeyBytes := btcecPublicKey.SerializeCompressed()

	pubkeyHash := btcutil.Hash160(pubkeyBytes)
	pubScript := append([]byte{0x00, 0x14}, pubkeyHash...)

	addr, err := btcutil.NewAddressScriptHash(pubScript, net)
	if err != nil {
		fmt.Println("NewAddressScriptHash failed ", err)
		return
	}
	fmt.Println("wif ", wif.String())
	fmt.Println("address ", addr.EncodeAddress())
	fmt.Println("expect  ", expectAddr)
	fmt.Println("pubScript ", hex.EncodeToString(pubScript))

}
