package utils

import (
	"bytes"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	MAX_LEAFS = 1 << 9
)

func Hash(data []byte) []byte {
	return crypto.Keccak256(data)
}

func MergeNodes(a, b []byte) []byte {
	var combined []byte
	if bytes.Compare(a, b) < 0 {
		combined = append(a, b...)
	} else {
		combined = append(b, a...)
	}
	return crypto.Keccak256(combined)
}

func Verify(proof [][]byte, root []byte, data []byte) bool {
	hashedLeaf := crypto.Keccak256(data)
	currentHash := hashedLeaf

	for _, p := range proof {
		if bytes.Compare(currentHash, p) < 0 {
			currentHash = crypto.Keccak256(append(currentHash, p...))
		} else {
			currentHash = crypto.Keccak256(append(p, currentHash...))
		}
	}

	return bytes.Equal(currentHash, root)
}
