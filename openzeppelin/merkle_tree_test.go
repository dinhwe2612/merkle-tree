package openzeppelin

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

var TestCases = []struct {
	name       string
	numLeaves  int
	numQueries int
}{
	{"Small", 1000, 1000},
	{"Medium", 100000, 100000},
	{"Large", 5000000, 5000000},
}

func Validate32Bytes(data []byte) bool {
	if len(data) != 32 {
		return false
	}
	for _, b := range data {
		if b < 0 || b > 255 {
			return false
		}
	}
	return true
}

func ValidateProof(proof [][]byte) bool {
	for _, p := range proof {
		if Validate32Bytes(p) == false {
			return false
		}
	}
	return true
}

func TestMerkleTree(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Logf("\n--- STRESS TEST: Fixed Merkle Tree with %d leaves, %d queries ---", tc.numLeaves, tc.numQueries)

			tree := &MerkleTree{}
			tree.Init(tc.numLeaves)

			leaves := make([]string, tc.numLeaves)
			for i := 0; i < tc.numLeaves; i++ {
				leaves[i] = fmt.Sprintf("leaf_%d", i)
			}

			startAdd := time.Now()
			for _, v := range leaves {
				tree.AddLeaf(v)
			}
			addTime := time.Since(startAdd)
			t.Logf("Added %d leaves in: %v (average %.3f µs/leaf)",
				tc.numLeaves, addTime, float64(addTime.Microseconds())/float64(tc.numLeaves))

			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			totalProof := time.Duration(0)
			totalVerify := time.Duration(0)
			okCount := 0

			// Successful queries
			numSuccessful := int(float64(tc.numQueries) * 0.8)
			for i := 0; i < numSuccessful; i++ {
				data := leaves[rnd.Intn(tc.numLeaves)]
				leaf := [32]byte(Keccak256([]byte(data)))

				startProof := time.Now()
				proof := tree.GetProof(data)
				proofTime := time.Since(startProof)

				startVerify := time.Now()
				ok := Verify(proof, tree.GetRoot(), leaf)
				verifyTime := time.Since(startVerify)

				if !ok {
					t.Fatalf("FAIL: Proof verification failed for leaf '%s' at query %d", leaf, i+1)
				}
				totalProof += proofTime
				totalVerify += verifyTime
				okCount++
			}

			// Failed queries
			numFailed := tc.numQueries - numSuccessful
			for i := 0; i < numFailed; i++ {
				leaf := fmt.Sprintf("leaf_%d", tc.numLeaves+rnd.Int())

				startProof := time.Now()
				proof := tree.GetProof(leaf)
				proofTime := time.Since(startProof)

				if proof != nil {
					t.Fatalf("FAIL: Expected no proof for non-existent leaf '%s' at query %d", leaf, i+1+numSuccessful)
				}
				totalProof += proofTime
				okCount++
			}

			t.Logf("%d queries: total proof time %v (average %.3f µs), total verification time %v (average %.3f µs), correct: %d/%d",
				tc.numQueries, totalProof, float64(totalProof.Microseconds())/float64(tc.numQueries),
				totalVerify, float64(totalVerify.Microseconds())/float64(tc.numQueries),
				okCount, tc.numQueries)
		})
	}
}
