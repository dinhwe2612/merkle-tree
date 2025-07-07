package merkletree

import (
	"bytes"
	"fmt"
	"sync"

	"merkle_module/utils"
)

type MerkleTree struct {
	nodes    [][]byte
	leafMap  map[string]int // map to store leaf hashes and their positions
	numLeafs int
	maxLeafs int
	treeID   int
	mu sync.Mutex // mutex to ensure thread safety
	    // mutex for thread safety
}

func NewMerkleTree(datas [][]byte, treeID int) (*MerkleTree, error) {
	tree := &MerkleTree{}
	tree.Init(utils.MAX_LEAFS)
	tree.build(datas)
	tree.treeID = treeID

	return tree, nil
}

func (tree *MerkleTree) Init(maxLeafs int) {
	if maxLeafs <= 0 {
		maxLeafs = utils.MAX_LEAFS
	}
	tree.maxLeafs = maxLeafs
	tree.nodes = make([][]byte, tree.maxLeafs<<1)
	tree.leafMap = make(map[string]int, tree.maxLeafs)
}

func (tree *MerkleTree) build(datas [][]byte) error {
	if len(datas) == 0 {
		// No data to build the tree
		return nil
	}

	// build the leaf map
	for i, data := range datas {
		tree.leafMap[string(data)] = i + 1 // store position starting from 1
		tree.nodes[tree.maxLeafs+i] = data
	}

	tree.numLeafs = len(datas)

	// compute hashes for parent nodes
	for parentStart := tree.maxLeafs >> 1; parentStart >= 1; parentStart >>= 1 {
		for nodeID := parentStart; nodeID < parentStart<<1; nodeID++ {
			leftChild := tree.nodes[nodeID<<1]
			rightChild := tree.nodes[nodeID<<1|1]
			tree.nodes[nodeID] = utils.MergeNodes(leftChild, rightChild)
		}
	}

	return nil
}

func (tree *MerkleTree) AddLeaf(data []byte) {
	tree.mu.Lock()
	defer tree.mu.Unlock()
	tree.numLeafs++
	tree.leafMap[string(data)] = tree.numLeafs // store position starting from 1
	tree.Update(data, tree.leafMap[string(data)])
}

func (tree *MerkleTree) Update(data []byte, pos int) {
	nodeID := tree.maxLeafs + pos - 1
	tree.nodes[nodeID] = data

	for nodeID > 1 {
		// println("Updating node:", nodeID, "with hash:", hash, "at position:", pos, "with parent:", (nodeID >> 1), "and sibling:", (nodeID ^ 1))
		parentID := nodeID >> 1
		siblingID := nodeID ^ 1
		tree.nodes[parentID] = utils.MergeNodes(tree.nodes[nodeID], tree.nodes[siblingID])
		nodeID = parentID
	}
}

func (tree *MerkleTree) GetMerkleRoot() []byte {
	if len(tree.leafMap) == 0 {
		return []byte{}
	}

	return tree.nodes[1]
}

func (tree *MerkleTree) GetProof(pos int) ([][]byte, error) {
	if pos <= 0 || pos > tree.numLeafs {
		return nil, fmt.Errorf("invalid position: %d, must be between 1 and %d", pos, tree.numLeafs)
	}

	proof := make([][]byte, 0, tree.maxLeafs)
	nodeID := tree.maxLeafs + pos - 1
	for nodeID > 1 {
		siblingID := nodeID ^ 1
		proof = append(proof, tree.nodes[siblingID])
		nodeID >>= 1
	}

	return proof, nil
}

func (tree *MerkleTree) GetListNodesToSave() []int {
	firstLeafID := tree.maxLeafs
	lastLeafID := firstLeafID + tree.numLeafs - 1
	nodesToSave := make([]int, 0, tree.numLeafs)
	for depth := 1; depth <= tree.maxLeafs; depth++ {
		for nodeID := firstLeafID; nodeID < lastLeafID; nodeID++ {
			if bytes.Compare(tree.nodes[nodeID], []byte{}) != 0 {
				nodesToSave = append(nodesToSave, nodeID)
			}
		}
		firstLeafID >>= 1
		lastLeafID >>= 1
	}

	return nodesToSave
}

func (tree *MerkleTree) GetTreeID() int {
	return tree.treeID
}

func (tree *MerkleTree) GetLastNodeID() (int, error) {
	if tree.numLeafs == 0 {
		return 0, fmt.Errorf("no nodes in the tree")
	}

	return tree.numLeafs, nil
}

func (tree *MerkleTree) Contains(data []byte) bool {
	_, exists := tree.leafMap[string(data)]
	return exists
}

func (tree *MerkleTree) IsFull() bool {
	return len(tree.leafMap) >= tree.maxLeafs
}
