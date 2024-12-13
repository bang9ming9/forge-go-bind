package main

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type Bind struct {
	ABI      interface{} `json:"abi"`
	ByteCode struct {
		Object string `json:"object"`
	} `json:"bytecode"`
	Sigs map[string]string `json:"methodIdentifiers"`
}

func (b Bind) Args(pkg, name string) (
	types []string,
	abis []string,
	bytecodes []string,
	fsigs []map[string]string,
	pKG string,
	lang bind.Lang,
	libs map[string]string,
	aliases map[string]string,
) {
	pKG = pkg
	lang = bind.LangGo
	libs = make(map[string]string)
	aliases = make(map[string]string)

	var abi string
	switch v := b.ABI.(type) {
	case string:
		abi = v
	default:
		bytes, _ := json.Marshal(v)
		abi = string(bytes)
	}

	types = []string{name}
	abis = []string{abi}

	bytecodes = []string{b.ByteCode.Object}
	fsigs = []map[string]string{b.Sigs}

	return
}
