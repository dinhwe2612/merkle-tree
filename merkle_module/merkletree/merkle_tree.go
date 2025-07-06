package merkletree

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"merkle_module/utils"

	"github.com/ethereum/go-ethereum/crypto"
)

type MerkleTree struct {
	merkleTree []string
	leafMap    map[string]int
	numLeafs   int
	maxLeafs   int
}

func NewMerkleTree(datas [][]byte) (*MerkleTree, error) {
	if len(datas) == 0 {
		return nil, nil
	}

	tree := &MerkleTree{}
	tree.Init(utils.MAX_LEAFS)
	tree.build(datas)

	return tree, nil
}

func (tree *MerkleTree) Init(maxLeafs int) {
	if maxLeafs <= 0 {
		maxLeafs = utils.MAX_LEAFS
	}
	tree.maxLeafs = maxLeafs
	tree.merkleTree = make([]string, tree.maxLeafs<<1)
	tree.leafMap = make(map[string]int, tree.maxLeafs)
}

func (tree *MerkleTree) build(datas [][]byte) error {
	if len(datas) == 0 {
		return fmt.Errorf("no datas provided to build the Merkle Tree")
	}

	// build leaf map
	for i, data := range datas {
		hash := utils.Hash(data)
		if _, exists := tree.leafMap[hash]; exists {
			return fmt.Errorf("duplicate leaf data found: %s", hash)
		}
		tree.leafMap[hash] = i + 1
		tree.merkleTree[tree.maxLeafs+i] = hash
	}

	tree.numLeafs = len(datas)

	// compute hashes for parent nodes
	for level := tree.maxLeafs >> 1; level >= 1; level >>= 1 {
		for nodeID := level; nodeID < level<<1; nodeID++ {
			leftChild := tree.merkleTree[nodeID<<1]
			rightChild := tree.merkleTree[nodeID<<1|1]
			tree.merkleTree[nodeID] = utils.MergeNodes(leftChild, rightChild)
		}
	}

	return nil
}

func (tree *MerkleTree) AddLeaf(data []byte) error {
	if len(tree.leafMap) >= tree.maxLeafs {
		return fmt.Errorf("Merkle Tree is full")
	}

	hash := utils.Hash(data)
	if _, exists := tree.leafMap[hash]; exists {
		return fmt.Errorf("leaf already exists")
	}

	tree.numLeafs++
	tree.leafMap[hash] = tree.numLeafs
	tree.Update(hash, tree.leafMap[hash])

	return nil
}

func (tree *MerkleTree) Update(hash string, pos int) {
	nodeID := tree.maxLeafs + pos - 1
	tree.merkleTree[nodeID] = hash

	for nodeID > 1 {
		// println("Updating node:", nodeID, "with hash:", hash, "at position:", pos, "with parent:", (nodeID >> 1), "and sibling:", (nodeID ^ 1))
		parentID := nodeID >> 1
		siblingID := nodeID ^ 1
		tree.merkleTree[parentID] = utils.MergeNodes(tree.merkleTree[nodeID], tree.merkleTree[siblingID])
		nodeID = parentID
	}
}

func (tree *MerkleTree) GetMerkleRoot() []byte {
	if len(tree.leafMap) == 0 {
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
	hash := utils.Hash(data)

	pos, exists := tree.leafMap[hash]
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

func (tree *MerkleTree) GetListNodesToSave() []int {
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
