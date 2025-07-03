package database

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	MAX_LEAFS = 1 << 15
)

type MerkleTree struct {
	merkleTree []string
	leafs      map[string]int
	numLeafs   int
	maxLeafs   int
	issuerDID  string
	id         string
}

func NewMerkleTree(data [][]byte) (*MerkleTree, error) {
	if len(data) == 0 {
		return nil, nil
	}
	tree := &MerkleTree{}
	tree.Init(MAX_LEAFS)
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
		tree.maxLeafs = MAX_LEAFS
	}
	tree.maxLeafs = maxLeafs
	tree.merkleTree = make([]string, tree.maxLeafs<<1)
	tree.leafs = make(map[string]int, tree.maxLeafs)
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
	tree.Update(hash, tree.leafs[hash])

	return nil
}

func (tree *MerkleTree) Update(hash string, pos int) {
	nodeID := tree.maxLeafs + pos - 1
	tree.merkleTree[nodeID] = hash

	for nodeID > 1 {
		// println("Updating node:", nodeID, "with hash:", hash, "at position:", pos, "with parent:", (nodeID >> 1), "and sibling:", (nodeID ^ 1))
		parentID := nodeID >> 1
		siblingID := nodeID ^ 1
		tree.merkleTree[parentID] = mergeNodes(tree.merkleTree[nodeID], tree.merkleTree[siblingID])
		nodeID = parentID
	}
}

func (tree *MerkleTree) GetMerkleRoot() []byte {
	if len(tree.leafs) == 0 {
		return []byte{}
	}

	rootNode := tree.merkleTree[1]

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

	proof := make([]string, 0, tree.maxLeafs)
	nodeID := tree.maxLeafs + pos - 1
	for nodeID > 1 {
		siblingID := nodeID ^ 1
		proof = append(proof, tree.merkleTree[siblingID])
		nodeID >>= 1
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

func (tree *MerkleTree) getListNodesToSave() []int {
	firstLeafID := tree.maxLeafs
	lastLeafID := firstLeafID + tree.numLeafs - 1
	nodesToSave := make([]int, 0, tree.numLeafs)
	for depth := 1; depth <= tree.maxLeafs; depth++ {
		for nodeID := firstLeafID; nodeID < lastLeafID; nodeID++ {
			if tree.merkleTree[nodeID] != "" {
				nodesToSave = append(nodesToSave, nodeID)
			}
		}
		firstLeafID >>= 1
		lastLeafID >>= 1
	}
	return nodesToSave
}

func (tree *MerkleTree) GetHashValues(nodeID int) (string, error) {
	if nodeID < 1 || nodeID >= len(tree.merkleTree) {
		return "", fmt.Errorf("node ID out of range")
	}
	if tree.merkleTree[nodeID] == "" {
		return "", fmt.Errorf("node does not exist")
	}
	return tree.merkleTree[nodeID], nil
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
