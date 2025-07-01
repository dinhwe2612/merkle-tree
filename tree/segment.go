package tree

const (
	MAX_SIZE = 1 << 20
)

type MerkleTree struct {
	merkleTree []string
	leafs      map[string]int
	numLeafs   int
	maxLeafs   int
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

func (tree *MerkleTree) AddLeaf(leaf string) {
	if len(tree.leafs) >= tree.maxLeafs {
		return
	}
	if _, exists := tree.leafs[leaf]; exists {
		// raise error if leaf already exists
		panic("Leaf already exists: " + leaf)
	}
	tree.numLeafs++
	tree.leafs[leaf] = tree.numLeafs
	tree.Update(leaf, tree.leafs[leaf], 1, 1, tree.maxLeafs)
}

func (tree *MerkleTree) Build(nodeID, begin, end int) {
	if begin == end {
		tree.merkleTree[nodeID] = Hash("#")
		return
	}
	mid := (begin + end) >> 1
	leftChild := nodeID << 1
	rightChild := nodeID<<1 | 1
	tree.Build(leftChild, begin, mid)
	tree.Build(rightChild, mid+1, end)
	tree.merkleTree[nodeID] = HashConcat(tree.merkleTree[leftChild], tree.merkleTree[rightChild])
}

func (tree *MerkleTree) Update(leaf string, pos int, nodeID, begin, end int) {
	if begin > end {
		return
	}
	if begin == end {
		tree.merkleTree[nodeID] = HashConcat(leaf, "")
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
	tree.merkleTree[nodeID] = HashConcat(tree.merkleTree[leftChild], tree.merkleTree[rightChild])
}

func (tree *MerkleTree) GetRoot() string {
	if len(tree.leafs) == 0 {
		return ""
	}
	return tree.merkleTree[1]
}

func (tree *MerkleTree) GetProof(leaf string) []ProofStep {
	pos, exists := tree.leafs[leaf]
	if !exists || pos == -1 {
		return nil
	}
	proof := make([]ProofStep, 0)
	nodeID, begin, end := 1, 1, tree.maxLeafs
	for begin < end {
		mid := (begin + end) >> 1
		if pos <= mid {
			proof = append(proof, ProofStep{key: tree.merkleTree[nodeID<<1|1], isLeft: false})
			end = mid
			nodeID <<= 1
		} else {
			proof = append(proof, ProofStep{key: tree.merkleTree[nodeID<<1], isLeft: true})
			begin = mid + 1
			nodeID = nodeID<<1 | 1
		}
	}
	for i, j := 0, len(proof)-1; i < j; i, j = i+1, j-1 {
		proof[i], proof[j] = proof[j], proof[i]
	}
	return proof
}

func prove(proof []ProofStep, leaf, merkleRoot string) bool {
	h := Hash(leaf)
	for _, n := range proof {
		if n.isLeft {
			h = HashConcat(n.key, h)
		} else {
			h = HashConcat(h, n.key)
		}
	}
	return h == merkleRoot
}
