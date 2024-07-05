package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func P2WPKHExample() {
	priv := []byte("btc-example-data")

	networkParams := &chaincfg.MainNetParams

	btcecPriKey, btcecPublicKey := btcec.PrivKeyFromBytes(priv)

	pubkeyHash := btcutil.Hash160(btcecPublicKey.SerializeCompressed())
	btcAddress, err := btcutil.NewAddressWitnessPubKeyHash(pubkeyHash, networkParams)
	if err != nil {
		fmt.Println("NewAddressWitnessPubKeyHash err ", err)
		return
	}

	wif, _ := btcutil.NewWIF(btcecPriKey, networkParams, true)
	fmt.Printf("priv %v, address %v\n", wif.String(), btcAddress.EncodeAddress())

	pubkeyBytes, _ := hex.DecodeString("0279BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798")
	pubkeyHash1 := btcutil.Hash160(pubkeyBytes)
	btcAddress1, err := btcutil.NewAddressWitnessPubKeyHash(pubkeyHash1, networkParams)
	if err != nil {
		fmt.Println("NewAddressWitnessPubKeyHash err ", err)
		return
	}
	fmt.Printf("address %v\n expect %v\n", btcAddress1.EncodeAddress(), "bc1qw508d6qejxtdg4y5r3zarvary0c5xw7kv8f3t4")

	spendAddressStr := btcAddress.EncodeAddress()
	spendAddr, err := btcutil.DecodeAddress(spendAddressStr, networkParams)
	if err != nil {
		fmt.Println("DecodeAddress spendAddr err", err)
		return
	}

	spenderAddrByte, err := txscript.PayToAddrScript(spendAddr)
	if err != nil {
		log.Println("spendAddr PayToAddrScript err", err)
		return
	}
	fmt.Printf("spenderAddrByte %x\n", spenderAddrByte)
	hash := "1234567890adcdef1234567890abcdef1234567890adcdef1234567890abcdef"
	chainHashFrom, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		fmt.Println("invalid hash")
		return
	}
	outAddrStr := spendAddressStr
	outAddress, err := btcutil.DecodeAddress(outAddrStr, networkParams)
	if err != nil {
		log.Println("DecodeAddress outAddress err", err)
		return
	}

	outScript, err := txscript.PayToAddrScript(outAddress)
	if err != nil {
		log.Println("outAddress PayToAddrScript err", err)
		return
	}

	txOuts := []*wire.TxOut{{
		PkScript: outScript,
		Value:    90000,
	}}

	outPoint := &wire.OutPoint{
		Hash:  *chainHashFrom,
		Index: 0,
	}
	tx := &wire.MsgTx{
		TxIn: []*wire.TxIn{{
			PreviousOutPoint: *outPoint,
			//SignatureScript: script
			Sequence: 0xffffffff,
		}},
		TxOut: txOuts,
	}

	tx.Version = 1
	tx.LockTime = 0x12345678

	a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		*outPoint: &wire.TxOut{
			Value:    100000,
			PkScript: spenderAddrByte,
		},
	})
	sigHashes := txscript.NewTxSigHashes(tx, a)

	sig, err := txscript.WitnessSignature(tx, sigHashes, 0, 100000, spenderAddrByte, txscript.SigHashDefault, wif.PrivKey, true)
	if err != nil {
		fmt.Println("WitnessSignature err ", err)
		return
	}
	fmt.Printf("sig %v\n", sig)

	tx.TxIn[0].Witness = sig

	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println(hexSignedTx)

}
