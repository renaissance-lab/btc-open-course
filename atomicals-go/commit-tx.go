package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func GetIntermediaAddressAndScript(wifPrivKey string, networkParam *chaincfg.Params) (string, string, error) {

	return GenerateP2TRAtomScriptAddress(wifPrivKey, "atom", "dmt", "0000", "", "sophon", 10000, 1000, networkParam)

}
func CommitExampleTx() {

	priv := []byte("btc-example-data-atom-key")
	chainParam := &chaincfg.TestNet3Params
	btcecPriKey, _ := btcec.PrivKeyFromBytes(priv)
	wif, err := btcutil.NewWIF(btcecPriKey, chainParam, true)
	if err != nil {
		fmt.Println("NewWIF err", err)
		return
	}

	spendAddrStr, err := GenerateP2TRAddress(wif.String(), chainParam)
	if err != nil {
		fmt.Println("GenerateP2TRAddress err ", err)
		return
	}
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

	chainHashFrom, err := chainhash.NewHashFromStr("157a9eefccfb0cf4740fbabc727e07faf57d1b5d75bdb20f72019cb65fc0d39c")
	if err != nil {
		fmt.Println("invalid hash, err ", err)
		return
	}
	outAddrStr, scriptData, err := GetIntermediaAddressAndScript(wif.String(), chainParam)
	if err != nil {
		fmt.Println("GetIntermediaAddressAndScript err ", err)
		return
	}
	fmt.Println("IntermediaAddress script data ", scriptData)
	outAddress, err := btcutil.DecodeAddress(outAddrStr, chainParam)
	if err != nil {
		log.Println("DecodeAddress outAddress err", err)
		return
	}

	outScript, err := txscript.PayToAddrScript(outAddress)
	if err != nil {
		log.Println("outAddress PayToAddrScript err", err)
		return
	}

	var getSeq uint32
	for i := uint32(0); i < 4294967295; i++ {
		tx := &wire.MsgTx{
			TxIn: []*wire.TxIn{{
				PreviousOutPoint: wire.OutPoint{
					Hash:  *chainHashFrom,
					Index: 2,
				},
				//SignatureScript: script,
				Witness:  [][]byte{[]byte("witness11111")},
				Sequence: i,
			}},
			TxOut: []*wire.TxOut{{
				PkScript: outScript,
				Value:    1000,
			}},
		}
		tx.Version = 1
		txHash := tx.TxHash().String()
		if strings.HasPrefix(txHash, "0000") {
			fmt.Println("tx ", txHash)
			getSeq = i
			break
		}
	}
	if getSeq > 0 { // get valid sequencer
		outPoint := &wire.OutPoint{
			Hash:  *chainHashFrom,
			Index: 2,
		}
		tx := &wire.MsgTx{
			TxIn: []*wire.TxIn{{
				PreviousOutPoint: *outPoint,
				//SignatureScript:  script,
				Sequence: getSeq,
			}},
			TxOut: []*wire.TxOut{{
				PkScript: outScript,
				Value:    1000,
			}},
		}
		tx.Version = 1

		a := txscript.NewMultiPrevOutFetcher(map[wire.OutPoint]*wire.TxOut{
			*outPoint: &wire.TxOut{
				Value:    2000,
				PkScript: spenderAddrByte,
			},
		})
		sigHashes := txscript.NewTxSigHashes(tx, a)

		signature, err := txscript.TaprootWitnessSignature(tx, sigHashes, 0, 2000, spenderAddrByte, txscript.SigHashDefault, wif.PrivKey)
		if err != nil {
			panic(err)
		}
		tx.TxIn[0].Witness = signature

		var signedTx bytes.Buffer
		tx.Serialize(&signedTx)

		hexSignedTx := hex.EncodeToString(signedTx.Bytes())
		fmt.Println(hexSignedTx)
	}
}
