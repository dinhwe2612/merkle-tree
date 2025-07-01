package openzeppelin

import (
	"bytes"
	"sort"

	"golang.org/x/crypto/sha3"
)

const (
	MAX_SIZE = 1 << 20
)

func Keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func mergeNodes(a, b [32]byte) [32]byte {
	// Create a slice of byte slices with actual data
	nodes := [][]byte{a[:], b[:]}

	// Sort the slices lexicographically
	sort.Slice(nodes, func(i, j int) bool {
		return bytes.Compare(nodes[i], nodes[j]) < 0
	})

	// Concatenate the sorted slices
	concatenated := append(nodes[0], nodes[1]...)

	// Compute Keccak256 hash (Ethereum uses legacy Keccak, not standard SHA3)
	hash := sha3.NewLegacyKeccak256()
	hash.Write(concatenated)

	var result [32]byte
	copy(result[:], hash.Sum(nil))
	return result
}

type MerkleTree struct {
	merkleTree [][32]byte
	leafs      map[[32]byte]int
	numLeafs   int
	maxLeafs   int
}

func (tree *MerkleTree) Init(maxLeafs int) {
	if maxLeafs <= 0 {
		tree.maxLeafs = MAX_SIZE
	}
	tree.maxLeafs = maxLeafs
	tree.merkleTree = make([][32]byte, tree.maxLeafs<<2)
	tree.leafs = make(map[[32]byte]int, tree.maxLeafs)
	tree.Build(1, 1, tree.maxLeafs)
	tree.numLeafs = 0
}

func (tree *MerkleTree) AddLeaf(data string) {
	if len(tree.leafs) >= tree.maxLeafs {
		return
	}
	leaf := [32]byte(Keccak256([]byte(data)))
	if _, exists := tree.leafs[leaf]; exists {
		// raise error if leaf already exists
		panic("Leaf already exists in the Merkle Tree")
	}
	tree.numLeafs++
	tree.leafs[leaf] = tree.numLeafs
	tree.Update(leaf, tree.leafs[leaf], 1, 1, tree.maxLeafs)
}

func (tree *MerkleTree) Build(nodeID, begin, end int) {
	if begin == end {
		tree.merkleTree[nodeID] = [32]byte(Keccak256([]byte("#")))
		return
	}
	mid := (begin + end) >> 1
	leftChild := nodeID << 1
	rightChild := nodeID<<1 | 1
	tree.Build(leftChild, begin, mid)
	tree.Build(rightChild, mid+1, end)
	tree.merkleTree[nodeID] = mergeNodes(tree.merkleTree[leftChild], tree.merkleTree[rightChild])
}

func (tree *MerkleTree) Update(leaf [32]byte, pos int, nodeID, begin, end int) {
	if begin > end {
		return
	}
	if begin == end {
		tree.merkleTree[nodeID] = [32]byte(Keccak256(leaf[:]))
		return
	}
	mid := (begin + end) >> 1
	leftChild := nodeID << 1
	rightChild := nodeID<<1 | 1
	if pos <= mid {
		tree.Update(leaf, pos, leftChild, begin, mid)
	} else {
		tree.Update(leaf, pos, rightChild, mid+1, end)
	}
	tree.merkleTree[nodeID] = mergeNodes(tree.merkleTree[leftChild], tree.merkleTree[rightChild])
}

func (tree *MerkleTree) GetRoot() [32]byte {
	if len(tree.leafs) == 0 {
		return [32]byte{}
	}
	return tree.merkleTree[1]
}

func (tree *MerkleTree) GetProof(data string) [][32]byte {
	leaf := [32]byte(Keccak256([]byte(data)))
	pos, exists := tree.leafs[leaf]
	if !exists || pos == -1 {
		return nil
	}
	proof := make([][32]byte, 0)
	nodeID := 1
	begin, end := 1, tree.maxLeafs
	for begin < end {
		mid := (begin + end) >> 1
		leftChild := nodeID << 1
		rightChild := nodeID<<1 | 1
		if pos <= mid {
			proof = append([][32]byte{[32]byte(Keccak256(tree.merkleTree[rightChild][:]))}, proof...)
			nodeID = leftChild
			end = mid
		} else {
			proof = append(proof, [32]byte(Keccak256(tree.merkleTree[leftChild][:])))
			nodeID = rightChild
			begin = mid + 1
		}
	}
	return proof
}

func Verify(proof [][32]byte, root [32]byte, leaf [32]byte) bool {
	if len(proof) == 0 {
		return bytes.Equal(root[:], leaf[:])
	}
	currentHash := leaf
	for _, p := range proof {
		if bytes.Compare(currentHash[:], p[:]) < 0 {
			currentHash = mergeNodes(currentHash, p)
		} else {
			currentHash = mergeNodes(p, currentHash)
		}
	}
	return bytes.Equal(currentHash[:], root[:])
}
