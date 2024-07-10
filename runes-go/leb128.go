package main

import (
	"encoding/hex"
	"fmt"
)

func encodeLEB128(n uint64) []byte {
	var buf []byte
	for {
		b := byte(n & 0x7f)
		n >>= 7
		if n == 0 {
			buf = append(buf, b)
			break
		}
		buf = append(buf, b|0x80)
	}
	return buf
}

func decodeLEB128(buf []byte) (uint64, int) {
	var n uint64
	var shift uint
	var i int
	for {
		b := buf[i]
		n |= uint64(b&0x7f) << shift
		shift += 7
		i++
		if b&0x80 == 0 {
			break
		}
	}
	return n, i
}

func EncodeExample() {
	fmt.Printf("encode 840566 %v\n", hex.EncodeToString(encodeLEB128(840566)))
	fmt.Printf("encode 2774 %v\n", hex.EncodeToString(encodeLEB128(2774)))
	byteData, _ := hex.DecodeString("b518")
	num, _ := decodeLEB128(byteData)
	fmt.Printf("decode b518 %v\n", num)
}
