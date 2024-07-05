package main

import (
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func P2TRExample1() {

	priv := []byte("btc-example-data")

	networkParams := &chaincfg.MainNetParams

	btcecPriKey, btcecPublicKey := btcec.PrivKeyFromBytes(priv)

	wif, _ := btcutil.NewWIF(btcecPriKey, networkParams, true)

	pubkey := btcecPublicKey

	legacyAddress, err := btcutil.NewAddressPubKey(wif.SerializePubKey(), networkParams)
	if err != nil {
		fmt.Println("NewAddressPubKey err ", err)
		return
	}
	fmt.Println("legacyAddress ", legacyAddress.EncodeAddress())

	taprootKey := txscript.ComputeTaprootKeyNoScript(pubkey)
	tapScriptAddr, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(taprootKey), networkParams,
	)
	if err != nil {
		fmt.Println("NewAddressTaproot err ", err)
		return
	}
	// address, err := btcutil.NewAddressTaproot(wif.SerializePubKey(), &chaincfg.MainNetParams)
	// if err != nil {
	// 	fmt.Println("NewAddressTaproot err ", err)
	// 	return "", err
	// }
	fmt.Println("address ", tapScriptAddr.String())

}
