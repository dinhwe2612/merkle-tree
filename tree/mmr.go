package tree

import (
	"strconv"
	"strings"
)

type MMRNode struct {
	left   *MMRNode
	right  *MMRNode
	parent *MMRNode
	hash   string
	height int
}

func mergeMMR(left, right *MMRNode) *MMRNode {
	newNode := &MMRNode{
		left:   left,
		right:  right,
		hash:   HashConcat(left.hash, right.hash),
		height: left.height + 1,
	}
	left.parent = newNode
	right.parent = newNode
	return newNode
}

type MMR struct {
	peaks     []*MMRNode
	leafs     []*MMRNode
	leafIndex map[string]int
}

func (tree *MMR) AddLeaf(data string) {
	h := Hash(data)
	newNode := &MMRNode{hash: h, height: 1}
	tree.leafs = append(tree.leafs, newNode)
	if tree.leafIndex == nil {
		tree.leafIndex = make(map[string]int)
	}
	tree.leafIndex[data] = len(tree.leafs) - 1
	for len(tree.peaks) > 0 && tree.peaks[len(tree.peaks)-1].height == newNode.height {
		last := tree.peaks[len(tree.peaks)-1]
		tree.peaks = tree.peaks[:len(tree.peaks)-1]
		newNode = mergeMMR(last, newNode)
	}
	tree.peaks = append(tree.peaks, newNode)
}

func (tree *MMR) GetRoot() string {
	var root strings.Builder
	root.WriteString(strconv.Itoa(len(tree.peaks)))
	for _, peak := range tree.peaks {
		root.WriteString(peak.hash)
	}
	return Hash(root.String())
}

func (tree *MMR) peakIndexOfLeafByValue(leaf string) int {
	idx, ok := tree.leafIndex[leaf]
	if !ok {
		return -1
	}
	count := 0
	for i, peak := range tree.peaks {
		size := 1 << (peak.height - 1)
		if idx < count+size {
			return i
		}
		count += size
	}
	return -1
}

func (tree *MMR) GetProof(leaf string) ([]ProofStep, string, string) {
	idx, ok := tree.leafIndex[leaf]
	if !ok {
		return nil, "", ""
	}
	peakIdx := tree.peakIndexOfLeafByValue(leaf)
	if peakIdx == -1 {
		return nil, "", ""
	}
	var peakPath []ProofStep
	node := tree.leafs[idx]
	for node.parent != nil {
		if node.parent.left == node {
			peakPath = append(peakPath, ProofStep{key: node.parent.right.hash, isLeft: false})
		} else {
			peakPath = append(peakPath, ProofStep{key: node.parent.left.hash, isLeft: true})
		}
		node = node.parent
	}
	var leftHash, rightHash strings.Builder
	leftHash.WriteString(strconv.Itoa(len(tree.peaks)))
	for i, peak := range tree.peaks {
		if i < peakIdx {
			leftHash.WriteString(peak.hash)
		} else if i > peakIdx {
			rightHash.WriteString(peak.hash)
		}
	}
	return peakPath, leftHash.String(), rightHash.String()
}

func VerifyProofMMR(leaf string, root string, peakPath []ProofStep, leftPeaks, rightPeaks string) bool {
	h := Hash(leaf)
	for _, step := range peakPath {
		if step.isLeft {
			h = Hash(step.key + h)
		} else {
			h = Hash(h + step.key)
		}
	}
	rootCalc := Hash(leftPeaks + h + rightPeaks)
	return rootCalc == root
}
