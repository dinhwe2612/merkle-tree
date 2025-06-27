package tree

import (
	"fmt"
	"math/rand"
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
	root := ""
	for _, peak := range mmr.peaks {
		root += peak.hash
	}
	return Hash(root)
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

type PeakProofStep struct {
	hash   string
	isLeft bool
}

type OtherPeakProofStep struct {
	hash      string
	peakIndex int
}

func (mmr *MMR) GetProofByValue(leaf string) ([]PeakProofStep, []OtherPeakProofStep) {
	idx, ok := mmr.leafIndex[leaf]
	if !ok {
		return nil, nil
	}
	peakIdx := mmr.peakIndexOfLeafByValue(leaf)
	if peakIdx == -1 {
		return nil, nil
	}
	var peakPath []PeakProofStep
	node := mmr.leafs[idx]
	for node.parent != nil {
		if node.parent.left == node {
			peakPath = append(peakPath, PeakProofStep{hash: node.parent.right.hash, isLeft: false})
		} else {
			peakPath = append(peakPath, PeakProofStep{hash: node.parent.left.hash, isLeft: true})
		}
		node = node.parent
	}
	var otherPeaks []OtherPeakProofStep
	for i, peak := range mmr.peaks {
		if i != peakIdx {
			otherPeaks = append(otherPeaks, OtherPeakProofStep{hash: peak.hash, peakIndex: i})
		}
	}
	return peakPath, otherPeaks
}

func VerifyProof(leaf string, root string, peakPath []PeakProofStep, otherPeaks []OtherPeakProofStep, leafPeakIndex int, totalPeaks int) bool {
	h := Hash(leaf)
	peaks := make([]string, totalPeaks)
	peaks[leafPeakIndex] = h
	for _, step := range peakPath {
		if step.isLeft {
			h = Hash(step.hash + h)
		} else {
			h = Hash(h + step.hash)
		}
	}
	peaks[leafPeakIndex] = h
	for _, step := range otherPeaks {
		peaks[step.peakIndex] = step.hash
	}
	var rootBuilder strings.Builder
	for _, p := range peaks {
		rootBuilder.WriteString(p)
	}
	rootCalc := Hash(rootBuilder.String())
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
		peakPath, otherPeaks := mmr.GetProofByValue(leaf)
		peakIdx := mmr.peakIndexOfLeafByValue(leaf)
		totalPeaks := len(mmr.peaks)
		proofTime := time.Since(startProof)
		startVerify := time.Now()
		ok := VerifyProof(leaf, mmr.Root(), peakPath, otherPeaks, peakIdx, totalPeaks)
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
