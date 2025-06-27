package tree

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type MMRNode struct {
	left   *MMRNode
	right  *MMRNode
	parent *MMRNode
	hash   string
	height int
}

func mergeMMR(left, right *MMRNode) *MMRNode {
	if left == nil {
		return right
	}
	if right == nil {
		return left
	}
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

func (mmr *MMR) AddLeaf(data string) {
	h := Hash(data)
	newNode := &MMRNode{hash: h, height: 1}
	mmr.leafs = append(mmr.leafs, newNode)
	if mmr.leafIndex == nil {
		mmr.leafIndex = make(map[string]int)
	}
	mmr.leafIndex[data] = len(mmr.leafs) - 1
	for len(mmr.peaks) > 0 && mmr.peaks[len(mmr.peaks)-1].height == newNode.height {
		last := mmr.peaks[len(mmr.peaks)-1]
		mmr.peaks = mmr.peaks[:len(mmr.peaks)-1]
		newNode = mergeMMR(last, newNode)
	}
	mmr.peaks = append(mmr.peaks, newNode)
}

func (mmr *MMR) Root() string {
	var root strings.Builder
	root.WriteString(strconv.Itoa(len(mmr.peaks)))
	for _, peak := range mmr.peaks {
		root.WriteString(peak.hash)
	}
	return Hash(root.String())
}

func (mmr *MMR) peakIndexOfLeafByValue(leaf string) int {
	idx, ok := mmr.leafIndex[leaf]
	if !ok {
		return -1
	}
	count := 0
	for i, peak := range mmr.peaks {
		size := 1 << (peak.height - 1)
		if idx < count+size {
			return i
		}
		count += size
	}
	return -1
}

func (mmr *MMR) GetProofByValue(leaf string) ([]ProofStep, string, string) {
	idx, ok := mmr.leafIndex[leaf]
	if !ok {
		return nil, "", ""
	}
	peakIdx := mmr.peakIndexOfLeafByValue(leaf)
	if peakIdx == -1 {
		return nil, "", ""
	}
	var peakPath []ProofStep
	node := mmr.leafs[idx]
	for node.parent != nil {
		if node.parent.left == node {
			peakPath = append(peakPath, ProofStep{key: node.parent.right.hash, isLeft: false})
		} else {
			peakPath = append(peakPath, ProofStep{key: node.parent.left.hash, isLeft: true})
		}
		node = node.parent
	}
	var leftHash, rightHash strings.Builder
	leftHash.WriteString(strconv.Itoa(len(mmr.peaks)))
	for i, peak := range mmr.peaks {
		if i < peakIdx {
			leftHash.WriteString(peak.hash)
		} else if i > peakIdx {
			rightHash.WriteString(peak.hash)
		}
	}
	return peakPath, leftHash.String(), rightHash.String()
}

func VerifyProof(leaf string, root string, peakPath []ProofStep, leftPeaks, rightPeaks string) bool {
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

func StressTestMMR(numLeaves int, numQueries int) {
	fmt.Printf("\n--- STRESS TEST: Fixed Merkle Tree with %d leaves, %d queries ---\n", numLeaves, numQueries)

	mmr := &MMR{}

	leaves := make([]string, numLeaves)
	for i := 0; i < numLeaves; i++ {
		leaves[i] = fmt.Sprintf("leaf_%d", i)
	}

	startAdd := time.Now()
	for _, v := range leaves {
		mmr.AddLeaf(v)
	}
	addTime := time.Since(startAdd)
	fmt.Printf("Added %d leaves in: %v (average %.3f µs/leaf)\n", numLeaves, addTime, float64(addTime.Microseconds())/float64(numLeaves))

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalProof := time.Duration(0)
	totalVerify := time.Duration(0)
	okCount := 0
	for t := 0; t < numQueries; t++ {
		leaf := leaves[rnd.Intn(numLeaves)]
		startProof := time.Now()
		peakPath, leftPeaks, rightPeaks := mmr.GetProofByValue(leaf)
		proofTime := time.Since(startProof)
		startVerify := time.Now()
		ok := VerifyProof(leaf, mmr.Root(), peakPath, leftPeaks, rightPeaks)
		verifyTime := time.Since(startVerify)
		totalProof += proofTime
		totalVerify += verifyTime
		if !ok {
			fmt.Printf("FAIL: Proof verification failed for leaf '%s' at query %d\n", leaf, t+1)
			return
		}
		okCount++
	}

	fmt.Printf("%d queries: total proof time %v (average %.3f µs), total verification time %v (average %.3f µs), correct: %d/%d\n",
		numQueries, totalProof, float64(totalProof.Microseconds())/float64(numQueries),
		totalVerify, float64(totalVerify.Microseconds())/float64(numQueries), okCount, numQueries)
}
