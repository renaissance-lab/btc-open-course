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

func RevealTxExample() {

	priv := []byte("btc-example-data-atom-key")
	chainParam := &chaincfg.TestNet3Params
	btcecPriKey, _ := btcec.PrivKeyFromBytes(priv)
	wif, err := btcutil.NewWIF(btcecPriKey, chainParam, true)
	if err != nil {
		fmt.Println("NewWIF err", err)
		return
	}

	spendAddrStr, scriptData, err := GetIntermediaAddressAndScript(wif.String(), chainParam)
	if err != nil {
		fmt.Println("GetIntermediaAddressAndScript err ", err)
		return
	}
	fmt.Println("spendAddrStr ", spendAddrStr)
	spendAddr, err := btcutil.DecodeAddress(spendAddrStr, chainParam)
	if err != nil {
		log.Println("DecodeAddress spendAddr err", err)
		return
	}

	spenderAddrByte, err := txscript.PayToAddrScript(spendAddr)
	if err != nil {
		log.Println("spendAddr PayToAddrScript err", err)
		return
	}

	chainHashFrom, err := chainhash.NewHashFromStr("000046a135535fde1c2103bbc1b6b420ac1a1105b3c1548c284d2f3c23313f9d")
	if err != nil {
		fmt.Println("invalid hash")
		return
	}

	outAddrStr, err := GenerateP2TRAddress(wif.String(), chainParam)
	if err != nil {
		fmt.Println("GenerateP2TRAddress err ", err)
		return
	}
	outAddr, err := btcutil.DecodeAddress(outAddrStr, chainParam)
	if err != nil {
		log.Println("DecodeAddress outAddrStr err", err)
		return
	}

	outAddrByte, err := txscript.PayToAddrScript(outAddr)
	if err != nil {
		log.Println("outAddr PayToAddrScript err", err)
		return
	}

	pkScript, err := hex.DecodeString(scriptData)
	if err != nil {
		log.Println("pkScripterr", err)
		return
	}
	tapLeaf := txscript.NewBaseTapLeaf(pkScript)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaf)

	internalKey := wif.PrivKey.PubKey()
	ctrlBlock := tapScriptTree.LeafMerkleProofs[0].ToControlBlock(
		internalKey,
	)

	tapScriptRootHash := tapScriptTree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(
		internalKey, tapScriptRootHash[:],
	)
	p2trScript, err := txscript.PayToTaprootScript(outputKey)
	if err != nil {
		log.Println("p2trScript", err)
		return
	}
	fmt.Printf("p2trScript      %x\n", p2trScript)
	fmt.Printf("spenderAddrByte %x\n", spenderAddrByte)
	outPoint := &wire.OutPoint{
		Hash:  *chainHashFrom,
		Index: 0,
	}
	tx := &wire.MsgTx{
		TxIn: []*wire.TxIn{{
			PreviousOutPoint: *outPoint,
			//SignatureScript:  script,
			Sequence: 0xfffe,
		}},
		TxOut: []*wire.TxOut{{
			PkScript: outAddrByte,
			Value:    800,
		}},
	}
	tx.Version = 1

	a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
		*outPoint: &wire.TxOut{
			Value:    1000,
			PkScript: spenderAddrByte,
		},
	})
	sigHashes := txscript.NewTxSigHashes(tx, a)

	// Now that we have the sig, we'll make a valid witness
	// including the control block.

	signature, err := txscript.RawTxInTapscriptSignature(tx, sigHashes, 0, 1000, spenderAddrByte, tapLeaf, txscript.SigHashDefault, wif.PrivKey)
	if err != nil {
		panic(err)
	}
	ctrlBlockBytes, err := ctrlBlock.ToBytes()
	if err != nil {
		fmt.Println("ctrlBlock.ToBytes() err ", err)
		return
	}
	txCopy := tx.Copy()
	txCopy.TxIn[0].Witness = wire.TxWitness{
		signature, pkScript, ctrlBlockBytes,
	}
	// tx.TxIn[0].Witness = signature

	var signedTx bytes.Buffer
	txCopy.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())
	fmt.Println(hexSignedTx)

}
