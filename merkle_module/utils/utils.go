package utils

import (
	"bytes"
	"merkle_module/domain/entities"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	MAX_LEAFS = 1 << 5
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

func GetTreeIDs(trees []*entities.MerkleTree) []int {
	ids := make([]int, len(trees))
	for i, tree := range trees {
		ids[i] = tree.ID
	}
	return ids
}

func NodesToBytes(nodes []*entities.MerkleNode) [][]byte {
	data := make([][]byte, len(nodes))
	for i, node := range nodes {
		data[i] = node.Data
	}
	return data
}

func ToByte32(data []byte) [32]byte {
	var byte32 [32]byte
	copy(byte32[:], data)
	return byte32
}
