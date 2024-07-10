package main

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

func GenerateP2TRAddress(keyStr, expectAddress string, netParam *chaincfg.Params) (string, error) {

	wif, err := btcutil.DecodeWIF(keyStr)
	if err != nil {
		fmt.Printf("DecodeWIF %v, err %v ", keyStr, err)
		return "", err
	}

	pubkey := wif.PrivKey.PubKey()

	taprootKey := txscript.ComputeTaprootKeyNoScript(pubkey)

	tapScriptAddr, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(taprootKey), netParam,
	)
	if err != nil {
		return "", err
	}
	if tapScriptAddr.String() != expectAddress && expectAddress != "" {
		errInfo := fmt.Sprintf("tapScriptAddr %v, expect %v", tapScriptAddr.String(), expectAddress)
		panic(errInfo)
	}

	return tapScriptAddr.String(), nil
}
