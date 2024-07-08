package main

import (
	"lukechampine.com/uint128"
)

type RuneData struct {
}

type TermsData struct {
}

type Etching struct {
	Divisibility uint8
	Premine      *uint128.Uint128
	Rune         *RuneData
	Spacers      uint32
	Symbol       []byte
	Terms        *TermsData
}
