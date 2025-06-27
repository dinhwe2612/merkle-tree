package tree

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	MAX_SIZE = 1 << 24
)

type MerkleTree struct {
	merkleTree []string
	leafs      map[string]int
	cnt        int
}

func (mt *MerkleTree) Init() {
	mt.merkleTree = make([]string, MAX_SIZE<<2)
	mt.leafs = make(map[string]int, MAX_SIZE)
	mt.Build(1, 1, MAX_SIZE)
	mt.cnt = 0
}

func (mt *MerkleTree) AddLeaf(leaf string) {
	if len(mt.leafs) >= MAX_SIZE {
		return
	}
	if _, exists := mt.leafs[leaf]; exists {
		// raise error if leaf already exists
		panic("Leaf already exists: " + leaf)
	}
	mt.cnt++
	mt.leafs[leaf] = mt.cnt
	mt.Update(leaf, 1, mt.leafs[leaf], 1, MAX_SIZE)
}

func (mt *MerkleTree) Build(p, b, e int) {
	if b == e {
		mt.merkleTree[p] = Hash("#")
		return
	}
	mid := (b + e) >> 1
	mt.Build(p<<1, b, mid)
	mt.Build(p<<1|1, mid+1, e)
	mt.merkleTree[p] = HashConcat(mt.merkleTree[p<<1], mt.merkleTree[p<<1|1])
}

func (mt *MerkleTree) Update(leaf string, pos int, p, b, e int) {
	if p < b || p > e {
		return
	}
	if b == e {
		mt.merkleTree[pos] = HashConcat(leaf, "")
		return
	}
	mid := (b + e) >> 1
	if p <= mid {
		mt.Update(leaf, pos<<1, p, b, mid)
	} else {
		mt.Update(leaf, pos<<1|1, p, mid+1, e)
	}
	mt.merkleTree[pos] = HashConcat(mt.merkleTree[pos<<1], mt.merkleTree[pos<<1|1])
}

func (mt *MerkleTree) GetRoot() string {
	if len(mt.leafs) == 0 {
		return ""
	}
	return mt.merkleTree[1]
}

type ProofStep struct {
	key    string
	isLeft bool
}

func (mt *MerkleTree) GetProof(leaf string) []ProofStep {
	pos, exists := mt.leafs[leaf]
	if !exists || pos == -1 {
		return nil
	}
	proof := make([]ProofStep, 0)
	p, b, e := 1, 1, MAX_SIZE
	for b < e {
		mid := (b + e) >> 1
		if pos <= mid {
			proof = append(proof, ProofStep{key: mt.merkleTree[p<<1|1], isLeft: false})
			e = mid
			p <<= 1
		} else {
			proof = append(proof, ProofStep{key: mt.merkleTree[p<<1], isLeft: true})
			b = mid + 1
			p = p<<1 | 1
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

func StressTestFixedMerkleTree(numLeaves int, numQueries int) {
	fmt.Printf("\n--- STRESS TEST: Fixed Merkle Tree with %d leaves, %d queries ---\n", numLeaves, numQueries)
	mt := &MerkleTree{}
	mt.Init()
	leaves := make([]string, numLeaves)
	for i := 0; i < numLeaves; i++ {
		leaves[i] = fmt.Sprintf("leaf_%d", i)
	}
	startAdd := time.Now()
	for _, v := range leaves {
		mt.AddLeaf(v)
	}
	addTime := time.Since(startAdd)
	fmt.Printf("Added %d leaves in: %v (average %.3f µs/leaf)\n", numLeaves, addTime, float64(addTime.Microseconds())/float64(numLeaves))

	fmt.Printf("\n--- Stress Test: Querying Proofs ---\n")
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalProof := time.Duration(0)
	totalVerify := time.Duration(0)
	okCount := 0
	for t := 0; t < numQueries; t++ {
		leaf := leaves[rnd.Intn(numLeaves)]
		startProof := time.Now()
		proof := mt.GetProof(leaf)
		proofTime := time.Since(startProof)
		startVerify := time.Now()
		ok := prove(proof, leaf, mt.GetRoot())
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
