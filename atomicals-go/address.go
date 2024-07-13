package main

import (
	"encoding/hex"
	"fmt"

	btcec "github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/txscript"
)

const (
	AtomicalsProtocolEnvelopeID = "atom"
)

func GenerateP2TRAddress(keyStr string, net *chaincfg.Params) (string, error) {
	wif, err := btcutil.DecodeWIF(keyStr)
	if err != nil {
		fmt.Println("DecodeWIF err ", err)
		return "", err
	}

	pubkey := wif.PrivKey.PubKey()

	taprootKey := txscript.ComputeTaprootKeyNoScript(pubkey)
	tapScriptAddr, err := btcutil.NewAddressTaproot(
		schnorr.SerializePubKey(taprootKey), net,
	)

	fmt.Println("address ", tapScriptAddr.String())
	fmt.Printf("tapScriptAddrTestNet WitnessProgram %x\n", tapScriptAddr.WitnessProgram())
	return tapScriptAddr.String(), err
}
func AppendMintUpdatedRevealScript(protocol, opType string, internalKey *btcec.PublicKey, payLoad AtomPayLoad) ([]byte, []byte, error) {

	builder := txscript.NewScriptBuilder()
	pubkey := schnorr.SerializePubKey(internalKey)
	builder.AddData(pubkey)
	// ops = append(ops, byte(txscript.OP_CHECKSIG))
	// ops = append(ops, byte(txscript.OP_0))
	// ops = append(ops, byte(txscript.OP_IF))
	builder.AddOp(txscript.OP_CHECKSIG)
	builder.AddOp(txscript.OP_0)
	builder.AddOp(txscript.OP_IF)

	// ops = append(ops, []byte(AtomicalsProtocolEnvelopeID)...)
	builder.AddData([]byte(protocol))

	// // optype
	// ops = append(ops, []byte(opType)...)
	builder.AddData([]byte(opType))
	// // data
	// payloadData := ToCbor(payLoad)
	// ops = append(ops, []byte(payloadData)...)
	payloadData := PayloadToCbor(payLoad)
	builder.AddData([]byte(payloadData))

	// ops = append(ops, byte(txscript.OP_ENDIF))
	builder.AddOp(txscript.OP_ENDIF)

	pkScript, err := builder.Script()
	if err != nil {
		return nil, nil, err
	}
	fmt.Printf("pkScript %x\n", pkScript)
	tapLeaf := txscript.NewBaseTapLeaf(pkScript)
	tapScriptTree := txscript.AssembleTaprootScriptTree(tapLeaf)
	tapScriptRootHash := tapScriptTree.RootNode.TapHash()
	outputKey := txscript.ComputeTaprootOutputKey(
		internalKey, tapScriptRootHash[:],
	)
	p2trScript, err := txscript.PayToTaprootScript(outputKey)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("p2trScript len ", len(p2trScript))
	fmt.Printf("p2trScript %x\n ", (p2trScript))
	calcPubkey := schnorr.SerializePubKey(outputKey)
	fmt.Printf("calcPubkey %x\n ", (calcPubkey))
	address, err := btcutil.NewAddressTaproot(calcPubkey, &chaincfg.MainNetParams)
	fmt.Printf("address %v, err %v\n", address, err)
	return pkScript, calcPubkey, nil

}
func GenerateP2TRAtomScriptAddress(keyStr, protocol, opType, bitworkr, bitworkc, mintTicker string, timeUnix, nonce int64, net *chaincfg.Params) (string, string, error) {
	wif, err := btcutil.DecodeWIF(keyStr)
	if err != nil {
		fmt.Println("DecodeWIF err ", err)
		return "", "", err
	}

	pubkey := wif.PrivKey.PubKey()
	// fmt.Printf("pubkey %v\n", pubkey.EncodeToString())
	args := AtomArg{
		Time:       timeUnix,
		Nonce:      nonce,
		Bitworkc:   bitworkc,
		Bitworkr:   bitworkr,
		MintTicker: mintTicker,
	}

	return GetScriptP2TRAddress(protocol, opType, pubkey, args, net)
}

func GetScriptP2TRAddress(protocol, opType string, internalKey *btcec.PublicKey, args AtomArg, net *chaincfg.Params) (string, string, error) {
	atomPayload := AtomPayLoad{Args: args}
	script, p2trPubkey, err := AppendMintUpdatedRevealScript(protocol, opType, internalKey, atomPayload)
	if err != nil {
		return "", "", err
	}
	address, err := btcutil.NewAddressTaproot(p2trPubkey, net)
	if err != nil {
		return "", "", err
	}
	// fmt.Printf("address for net %v, is %v\n", net, address.EncodeAddress())
	return address.EncodeAddress(), hex.EncodeToString(script), nil
}
