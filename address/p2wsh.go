package main

import (
	"crypto/sha256"
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func P2WSHExample() {
	net := &chaincfg.MainNetParams
	priv1 := []byte("btc-example-data1")
	priv2 := []byte("btc-example-data2")
	btcecPriKey1, btcecPublicKey1 := btcec.PrivKeyFromBytes(priv1)
	btcecPriKey2, btcecPublicKey2 := btcec.PrivKeyFromBytes(priv2)

	wif1, _ := btcutil.NewWIF(btcecPriKey1, net, true)
	fmt.Printf("priv1 %v\n", wif1.String())

	wif2, _ := btcutil.NewWIF(btcecPriKey2, net, true)
	fmt.Printf("priv2 %v\n", wif2.String())

	keysLen := 2
	requiredKeys := 2
	addressHashKeys := make([]*btcutil.AddressPubKey, keysLen)

	var err error
	addressHashKeys[0], err = btcutil.NewAddressPubKey(btcecPublicKey1.SerializeCompressed(), net)
	if err != nil {
		fmt.Println("NewAddressPubKey1 failed ", err)
		return
	}

	addressHashKeys[1], err = btcutil.NewAddressPubKey(btcecPublicKey2.SerializeCompressed(), net)
	if err != nil {
		fmt.Println("NewAddressPubKey2 failed ", err)
		return
	}

	pubScript, err := txscript.MultiSigScript(addressHashKeys, requiredKeys)
	if err != nil {
		fmt.Println("MultiSigScript failed ", err)
		return
	}
	sha256Calc := sha256.New()
	sha256Calc.Write(pubScript)
	witnessProg := sha256Calc.Sum(nil)
	fmt.Printf("witnessProg %x, len %v\n", witnessProg, len(witnessProg))
	addr, err := btcutil.NewAddressWitnessScriptHash(witnessProg, net)
	if err != nil {
		fmt.Println("NewAddressScriptHash failed ", err)
		return
	}
	fmt.Println("address ", addr.EncodeAddress())
}
