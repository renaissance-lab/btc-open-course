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

const (
	DustAmount = int64(546)
)

// suppose one input and 3 outputs
func RuneMintTxVSize(size int) int64 {
	return int64(120) + int64(size-1)*32
}

func CreateOrdiMintTx(wifKey, spendAddrStr, hash string, index uint32, amount, mintTokenStashiAmount, fee int64, chainParam *chaincfg.Params, mintData []byte) (string, int64, error) {

	wif, err := btcutil.DecodeWIF(wifKey)
	if err != nil {
		fmt.Println("DecodeWIF err ", err)
		return "", 0, err
	}
	// check sepndAddr
	addrGen, err := GenerateP2TRAddress(wifKey, spendAddrStr, chainParam)
	if err != nil {
		fmt.Println("GenerateP2TRAddress err ", err)
		return "", 0, err
	}
	if addrGen != spendAddrStr {
		fmt.Println("GenerateP2TRAddress addrGen %v, spendAddrStr %v not match", addrGen, spendAddrStr)
		return "", 0, fmt.Errorf("address not match")
	}

	// spendAddrStr := "tb1px8ap86kjz2cesyd6sy4r6cmfhpxe2lhd0zvshatm73zlz4wrnpesqfpxes"
	spendAddr, err := btcutil.DecodeAddress(spendAddrStr, chainParam)
	if err != nil {
		log.Println("DecodeAddress spendAddr err", err)
		return "", 0, err
	}

	spenderAddrByte, err := txscript.PayToAddrScript(spendAddr)
	if err != nil {
		log.Println("spendAddr PayToAddrScript err", err)
		return "", 0, err
	}

	chainHashFrom, err := chainhash.NewHashFromStr(hash)
	if err != nil {
		fmt.Println("invalid hash")
		return "", 0, err
	}

	fmt.Printf("spenderAddrByte %x\n", spenderAddrByte)

	outAddress := spendAddr // send to from address

	outScript, err := txscript.PayToAddrScript(outAddress)
	if err != nil {
		log.Println("outAddress PayToAddrScript err", err)
		return "", 0, err
	}
	// UTXO 0 OP_RETURN

	// UTXO 0   OP_RETRUN 6a5d0814c0a23314161601 SATOSHI 01
	// scriptPubKey, _ := hex.DecodeString("6a5d0814c0a23314161601")
	txOuts := []*wire.TxOut{{
		PkScript: mintData,
		Value:    0,
	}}
	// UTXO 1 : my address
	txOuts = append(txOuts, &wire.TxOut{
		PkScript: outScript,
		Value:    mintTokenStashiAmount,
	})

	size := RuneMintTxVSize(2)
	// there are redeem utxo
	redeemAmount := amount - 3*mintTokenStashiAmount - size*fee
	if amount > mintTokenStashiAmount+DustAmount+size*fee {
		txOuts = append(txOuts, &wire.TxOut{
			PkScript: spenderAddrByte,
			Value:    redeemAmount,
		})
	}
	tx := &wire.MsgTx{
		TxIn: []*wire.TxIn{{
			PreviousOutPoint: wire.OutPoint{
				Hash:  *chainHashFrom,
				Index: index,
			},
			//SignatureScript: script,
			Witness:  [][]byte{[]byte("witness11111")},
			Sequence: 0,
		}},
		TxOut: txOuts,
	}

	tx.Version = 1

	var w bytes.Buffer
	err = tx.SerializeNoWitness(&w)
	if err != nil {
		fmt.Println("SerializeNoWitness err ", err)
		return "", 0, err
	}
	serializedBlock := w.Bytes()
	fmt.Printf("serializedBlock %x， len(serializedBlock) is %v\n", serializedBlock, len(serializedBlock))

	// 01000000 - 交易的版本号
	// 01 - 交易输入的数量
	// afb466816ddcf2003bcc64a73b4c3ce627d5af42dd1654b8c4e35e894db73ada - 前序交易的 id
	// 00000000 - 前序输出的索引号
	// 00 - 传统签名占位符
	// 00000000 - 序列号
	// 01 - 输出的数量
	// 80f0fa0200000000 - 输出的数额，50,000,000 聪 = 0.5 bitcoin 1600140e6a5ae16b91296707c28ef5d0836f04667bdae3 - 输出的锁定脚本
	// 00000000 - 交易的时间锁

	tx.TxIn[0].Sequence = 0x88888888
	previousOutPoint0 := &wire.OutPoint{
		Hash:  *chainHashFrom,
		Index: index,
	}

	a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		*previousOutPoint0: &wire.TxOut{
			Value:    amount,
			PkScript: spenderAddrByte,
		},
	})
	sigHashes := txscript.NewTxSigHashes(tx, a)

	signature, err := txscript.TaprootWitnessSignature(tx, sigHashes, 0, amount, spenderAddrByte, txscript.SigHashDefault, wif.PrivKey)
	if err != nil {
		fmt.Println("TaprootWitnessSignature err", err)
		return "", 0, err
	}
	tx.TxIn[0].Witness = signature

	var signedTx bytes.Buffer
	tx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println(hexSignedTx)
	return tx.TxHash().String(), redeemAmount, nil
}

func MintExample() {
	priv := []byte("btc-example-data-runes-key")
	chainParam := &chaincfg.TestNet3Params
	btcecPriKey, _ := btcec.PrivKeyFromBytes(priv)
	wif, err := btcutil.NewWIF(btcecPriKey, chainParam, true)
	if err != nil {
		fmt.Println("NewWIF err", err)
		return
	}

	spendAddrStr, err := GenerateP2TRAddress(wif.String(), "", chainParam)
	if err != nil {
		fmt.Println("GenerateP2TRAddress err ", err)
		return
	}
	fmt.Println("spendAddrStr ", spendAddrStr)

	testHash := "157a9eefccfb0cf4740fbabc727e07faf57d1b5d75bdb20f72019cb65fc0d39c"
	testAmount := int64(100000)
	mintSatoshi := int64(5000)
	fee := int64(10)

	height := uint64(2609649)
	index := uint64(946)
	bytesData, err := GetMintScript(height, index, 1, false)
	if err != nil {
		fmt.Println("GetMintScript err ", err)
		return
	}

	CreateOrdiMintTx(wif.String(), spendAddrStr, testHash, 0, testAmount, mintSatoshi, fee, chainParam, bytesData)

}
