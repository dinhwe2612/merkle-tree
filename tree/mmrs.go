package tree

const (
	MAX_MMR_SIZE = 1 << 6 // Maximum number of leaves in one MMR
)

type MMRs struct {
	trees []MMR          // list of MMRs (0-indexed)
	leafs map[string]int // leaf belong to i-th MMR
}

func (tree *MMRs) AddLeaf(data string) {
	// check if current MMR is full
	lastIndex := len(tree.trees) - 1
	if len(tree.trees[lastIndex].leafs) >= MAX_MMR_SIZE {
		// create a new MMR
		newMMR := MMR{}
		tree.trees = append(tree.trees, newMMR)
	}

	tree.trees[lastIndex].AddLeaf(data)

	h := Hash(data)
	if _, ok := tree.leafs[h]; !ok {
		tree.leafs[h] = lastIndex
	} else {
		// raise error if leaf already exists
		panic("Leaf already exists: " + data)
	}
}

func (tree *MMRs) GetRoot(data string) string {
	h := Hash(data)
	if idx, exists := tree.leafs[h]; exists {
		return tree.trees[idx].GetRoot()
	}
	return ""
}

func (tree *MMRs) GetProofByValue(data string) ([]ProofStep, string, string) {
	h := Hash(data)
	if idx, exists := tree.leafs[h]; exists {
		peakIdx := tree.trees[idx].peakIndexOfLeafByValue(data)
		if peakIdx == -1 {
			return nil, "", ""
		}
		proof, leftHash, rightHash := tree.trees[idx].GetProof(data)
		// println("Proof for leaf:", data, "in MMR index:", idx, "peak index:", peakIdx)
		return proof, leftHash, rightHash
	}
	return nil, "", ""
}
