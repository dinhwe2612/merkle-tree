package openzeppelin

import (
	"bytes"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/crypto"
)

func Hash(data []byte) string {
	return hex.EncodeToString(crypto.Keccak256(data))
}

func mergeNodes(a, b string) string {
	aBytes, _ := hex.DecodeString(a)
	bBytes, _ := hex.DecodeString(b)
	var combined []byte
	if bytes.Compare(aBytes, bBytes) < 0 {
		combined = append(aBytes, bBytes...)
	} else {
		combined = append(bBytes, aBytes...)
	}
	return hex.EncodeToString(crypto.Keccak256(combined))
}
