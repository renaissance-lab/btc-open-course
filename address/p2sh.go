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

func MergeMultiSign(rawSigns [][]byte, pkScript []byte) ([]byte, error) {
	builder := txscript.NewScriptBuilder()

	builder.AddOp(txscript.OP_FALSE)
	for _, sign := range rawSigns {
		builder.AddData(sign)
	}
	builder.AddData(pkScript)

	return builder.Script()
}

func MultisignExample() {

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

	addr, err := btcutil.NewAddressScriptHash(pubScript, net)
	if err != nil {
		fmt.Println("NewAddressScriptHash failed ", err)
		return
	}
	fmt.Println("address ", addr.EncodeAddress())
	fmt.Println("pubScript ", hex.EncodeToString(pubScript))
	//fmt.Printf("multi address %v\n", addr.EncodeAddress())

	// 签名的话
	spendAddr := addr
	spendAddressStr := spendAddr.EncodeAddress()
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
	outAddress, err := btcutil.DecodeAddress(outAddrStr, net)
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

	sig1, err := txscript.RawTxInSignature(tx, 0, spenderAddrByte, txscript.SigHashAll, wif1.PrivKey)
	if err != nil {
		fmt.Println("RawTxInSignature1 err ", err)
		return
	}
	sig2, err := txscript.RawTxInSignature(tx, 0, spenderAddrByte, txscript.SigHashAll, wif2.PrivKey)
	if err != nil {
		fmt.Println("RawTxInSignature2 err ", err)
		return
	}

	signScript, err := MergeMultiSign([][]byte{sig1, sig2}, pubScript)
	if err != nil {
		fmt.Println("MergeMultiSign err ", err)
		return
	}
	tx.TxIn[0].SignatureScript = signScript

	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println(hexSignedTx)
}
