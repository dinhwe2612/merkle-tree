package openzeppelin

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	MAX_SIZE = 1 << 20
)

type MerkleTree struct {
	merkleTree []string
	leafs      map[string]int
	numLeafs   int
	maxLeafs   int
}

func NewMerkleTree(data [][]byte) (*MerkleTree, error) {
	if len(data) == 0 {
		return nil, nil
	}
	tree := &MerkleTree{}
	tree.Init(MAX_SIZE)
	for _, item := range data {
		err := tree.AddLeaf(item)
		if err != nil {
			return nil, fmt.Errorf("failed to add leaf: %v", err)
		}
	}
	return tree, nil
}

func (tree *MerkleTree) Init(maxLeafs int) {
	if maxLeafs <= 0 {
		tree.maxLeafs = MAX_SIZE
	}
	tree.maxLeafs = maxLeafs
	tree.merkleTree = make([]string, tree.maxLeafs<<2)
	tree.leafs = make(map[string]int, tree.maxLeafs)
	tree.Build(1, 1, tree.maxLeafs)
	tree.numLeafs = 0
}

func (tree *MerkleTree) AddLeaf(data []byte) error {
	if len(tree.leafs) >= tree.maxLeafs {
		return fmt.Errorf("Merkle Tree is full")
	}
	hash := hex.EncodeToString(crypto.Keccak256(data))
	if _, exists := tree.leafs[hash]; exists {
		return fmt.Errorf("leaf already exists")
	}
	tree.numLeafs++
	tree.leafs[hash] = tree.numLeafs
	tree.Update(hash, tree.leafs[hash], 1, 1, tree.maxLeafs)
	return nil
}

func (tree *MerkleTree) Build(nodeID, begin, end int) {
	if begin == end {
		tree.merkleTree[nodeID] = hex.EncodeToString(crypto.Keccak256([]byte("#")))
		return
	}
	mid := (begin + end) >> 1
	leftChild := nodeID << 1
	rightChild := nodeID<<1 | 1
	tree.Build(leftChild, begin, mid)
	tree.Build(rightChild, mid+1, end)
	tree.merkleTree[nodeID] = mergeNodes(tree.merkleTree[leftChild], tree.merkleTree[rightChild])
}

func (tree *MerkleTree) Update(hash string, pos, nodeID, begin, end int) {
	if begin > end {
		return
	}
	if begin == end {
		tree.merkleTree[nodeID] = hash
		return
	}
	mid := (begin + end) >> 1
	leftChild := nodeID << 1
	rightChild := nodeID<<1 | 1
	if pos <= mid {
		tree.Update(hash, pos, leftChild, begin, mid)
	} else {
		tree.Update(hash, pos, rightChild, mid+1, end)
	}
	tree.merkleTree[nodeID] = mergeNodes(tree.merkleTree[leftChild], tree.merkleTree[rightChild])
}

func (tree *MerkleTree) GetMerkleRoot() []byte {
	if len(tree.leafs) == 0 {
		return []byte{}
	}
	rootNode := tree.merkleTree[1]
	if rootNode == "" {
		return []byte{}
	}
	rootBytes, err := hex.DecodeString(rootNode)
	if err != nil {
		fmt.Printf("Error decoding root node: %v\n", err)
		return []byte{}
	}
	return rootBytes
}

func (tree *MerkleTree) GetProof(data []byte) ([][]byte, error) {
	hash := hex.EncodeToString(crypto.Keccak256(data))
	pos, exists := tree.leafs[hash]
	if !exists {
		return nil, fmt.Errorf("leaf not found")
	}
	proof := []string{}
	nodeID := 1
	begin, end := 1, tree.maxLeafs
	for begin < end {
		mid := (begin + end) >> 1
		leftChild := nodeID << 1
		rightChild := nodeID<<1 | 1
		if pos <= mid {
			proof = append(proof, tree.merkleTree[rightChild])
			nodeID = leftChild
			end = mid
		} else {
			proof = append(proof, tree.merkleTree[leftChild])
			nodeID = rightChild
			begin = mid + 1
		}
	}
	// reverse the proof
	for i, j := 0, len(proof)-1; i < j; i, j = i+1, j-1 {
		proof[i], proof[j] = proof[j], proof[i]
	}
	returnBytes := make([][]byte, len(proof))
	for i, p := range proof {
		bytes, err := hex.DecodeString(p)
		if err != nil {
			return nil, fmt.Errorf("error decoding proof: %v", err)
		}
		returnBytes[i] = bytes
	}
	return returnBytes, nil
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
