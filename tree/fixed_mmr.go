package tree

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	MAX_MMR_SIZE = 1 << 6 // Maximum number of leaves in one MMR
)

type FixedMMRs struct {
	trees []MMR          // list of MMRs (0-indexed)
	leafs map[string]int // leaf belong to i-th MMR
}

func (fmmr *FixedMMRs) AddLeaf(data string) {
	// check if current MMR is full
	lastIndex := len(fmmr.trees) - 1
	if len(fmmr.trees[lastIndex].leafs) >= MAX_MMR_SIZE {
		// create a new MMR
		newMMR := MMR{}
		fmmr.trees = append(fmmr.trees, newMMR)
	}

	fmmr.trees[lastIndex].AddLeaf(data)

	h := Hash(data)
	if _, ok := fmmr.leafs[h]; !ok {
		fmmr.leafs[h] = lastIndex
	} else {
		// raise error if leaf already exists
		panic("Leaf already exists: " + data)
	}
}

func (fmmr *FixedMMRs) Root(data string) string {
	h := Hash(data)
	if idx, exists := fmmr.leafs[h]; exists {
		return fmmr.trees[idx].Root()
	}
	return ""
}

func (fmmr *FixedMMRs) GetProofByValue(data string) ([]ProofStep, string, string) {
	h := Hash(data)
	if idx, exists := fmmr.leafs[h]; exists {
		peakIdx := fmmr.trees[idx].peakIndexOfLeafByValue(data)
		if peakIdx == -1 {
			return nil, "", ""
		}
		proof, leftHash, rightHash := fmmr.trees[idx].GetProofByValue(data)
		// println("Proof for leaf:", data, "in MMR index:", idx, "peak index:", peakIdx)
		return proof, leftHash, rightHash
	}
	return nil, "", ""
}

func StressTestFixedMMRs(numLeaves int, numQueries int) {
	fmt.Printf("\n--- STRESS TEST: FixedMMRs with %d leaves, %d queries ---\n", numLeaves, numQueries)

	fmmr := &FixedMMRs{
		trees: []MMR{{}},
		leafs: make(map[string]int),
	}

	leaves := make([]string, numLeaves)
	for i := 0; i < numLeaves; i++ {
		leaves[i] = fmt.Sprintf("leaf_%d", i)
	}

	// Add leaves
	startAdd := time.Now()
	for _, leaf := range leaves {
		fmmr.AddLeaf(leaf)
	}
	addTime := time.Since(startAdd)
	fmt.Printf("Added %d leaves in: %v (average %.3f µs/leaf)\n",
		numLeaves, addTime, float64(addTime.Microseconds())/float64(numLeaves))

	// Proof and verification
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalProof := time.Duration(0)
	totalVerify := time.Duration(0)
	okCount := 0

	for i := 0; i < numQueries; i++ {
		leaf := leaves[rnd.Intn(numLeaves)]

		startProof := time.Now()
		proof, left, right := fmmr.GetProofByValue(leaf)
		proofTime := time.Since(startProof)

		startVerify := time.Now()
		ok := VerifyProof(leaf, fmmr.Root(leaf), proof, left, right)
		verifyTime := time.Since(startVerify)

		totalProof += proofTime
		totalVerify += verifyTime

		if !ok {
			fmt.Printf("FAIL: Proof verification failed for leaf '%s' at query %d\n", leaf, i+1)
			return
		}
		okCount++
	}

	fmt.Printf("%d queries: total proof time %v (avg %.3f µs), total verification time %v (avg %.3f µs), correct: %d/%d\n",
		numQueries, totalProof, float64(totalProof.Microseconds())/float64(numQueries),
		totalVerify, float64(totalVerify.Microseconds())/float64(numQueries), okCount, numQueries)
}
