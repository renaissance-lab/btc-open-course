package main

import (
	"encoding/hex"
	"fmt"

	"github.com/btcsuite/btcd/txscript"
)

const (
	MAGIC_NUMBER     = txscript.OP_13
	RUNE_MINT_TAG    = 0x14
	RUNE_POINTER_TAG = 0x16
)

func GenerateMintData(blockHeight, index uint64, outputIndex uint32, needSet bool) []byte {
	var retData []byte
	retData = append(retData, RUNE_MINT_TAG)
	heightBytes := encodeLEB128(blockHeight)
	retData = append(retData, heightBytes...)
	retData = append(retData, RUNE_MINT_TAG)
	indexBytes := encodeLEB128(index)
	retData = append(retData, indexBytes...)
	if needSet {
		retData = append(retData, RUNE_POINTER_TAG)
		outputIndexBytes := encodeLEB128(uint64(outputIndex))
		retData = append(retData, outputIndexBytes...)
	}
	return retData
}
func GetMintScript(blockHeight, index uint64, outputIndex uint32, needSet bool) ([]byte, error) {
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN)
	builder.AddOp(MAGIC_NUMBER)
	//
	chunkData := GenerateMintData(blockHeight, index, outputIndex, needSet)
	builder.AddData(chunkData)
	return builder.Script()
}

func MintScriptExample() {
	height := uint64(2609649)
	index := uint64(946)
	bytesData, err := GetMintScript(height, index, 1, false)
	if err != nil {
		fmt.Println("GetMintScript err ", err)
		return
	}
	fmt.Printf("Mint Rune[%v:%v] data 0x%v", height, index, hex.EncodeToString(bytesData))
}
