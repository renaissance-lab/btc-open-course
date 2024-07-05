package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	//btcec "github.com/btcsuite/btcd/btcec/v2"
	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func P2PKHExample() {

	priv := []byte("btc-example-data")

	networkParams := &chaincfg.MainNetParams

	btcecPriKey, btcecPublicKey := btcec.PrivKeyFromBytes(priv)
	btcAddress, err := btcutil.NewAddressPubKey(btcecPublicKey.SerializeCompressed(), networkParams)
	if err != nil {
		fmt.Println("NewAddressPubKey err ", err)
		return
	}

	wif, _ := btcutil.NewWIF(btcecPriKey, networkParams, true)
	fmt.Printf("priv %v, address %v\n", wif.String(), btcAddress.EncodeAddress())
	//return []byte(wif.String()), btcAddress.EncodeAddress(), nil
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

	tx := &wire.MsgTx{
		TxIn: []*wire.TxIn{{
			PreviousOutPoint: wire.OutPoint{
				Hash:  *chainHashFrom,
				Index: 0,
			},
			//SignatureScript: script
			Sequence: 0xffffffff,
		}},
		TxOut: txOuts,
	}

	tx.Version = 1
	tx.LockTime = 0x12345678

	sig, err := txscript.RawTxInSignature(tx, 0, spenderAddrByte, txscript.SigHashAll, wif.PrivKey)
	if err != nil {
		fmt.Println("RawTxInSignature err ", err)
		return
	}
	fmt.Printf("sig %v\n", hex.EncodeToString(sig))

	pubkeyBytes := btcecPublicKey.SerializeCompressed()
	fmt.Printf("pubkeyBytes %v\n", hex.EncodeToString(pubkeyBytes))

	signature, err := txscript.SignatureScript(tx, 0, spenderAddrByte, txscript.SigHashAll, wif.PrivKey, true)
	if err != nil {
		fmt.Println("RawTxInSignature err", err)
		return
	}
	fmt.Printf("signature  %x\n", signature)
	tx.TxIn[0].SignatureScript = signature

	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println(hexSignedTx)

}
