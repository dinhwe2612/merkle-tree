package tree

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestMMRs(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.name, func(t *testing.T) {
			// t.Logf("--- STRESS TEST: MMRs with %d leaves, %d queries ---", tc.numLeaves, tc.numQueries)

			tree := &MMRs{
				trees: []MMR{{}},
				leafs: make(map[string]int),
			}

			leaves := make([]string, tc.numLeaves)
			for i := range leaves {
				leaves[i] = fmt.Sprintf("leaf_%d", i)
			}

			startAdd := time.Now()
			for _, leaf := range leaves {
				tree.AddLeaf(leaf)
			}
			addTime := time.Since(startAdd)
			t.Logf("Added %d leaves in: %v (average %.3f µs/leaf)",
				tc.numLeaves, addTime, float64(addTime.Microseconds())/float64(tc.numLeaves))

			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			totalProof := time.Duration(0)
			totalVerify := time.Duration(0)
			okCount := 0

			// Successful queries
			numSuccess := int(float64(tc.numQueries) * 0.8)
			for i := 0; i < numSuccess; i++ {
				leaf := leaves[rnd.Intn(tc.numLeaves)]

				startProof := time.Now()
				proof, left, right := tree.GetProofByValue(leaf)
				proofTime := time.Since(startProof)

				startVerify := time.Now()
				ok := VerifyProofMMR(leaf, tree.GetRoot(leaf), proof, left, right)
				verifyTime := time.Since(startVerify)

				totalProof += proofTime
				totalVerify += verifyTime

				if !ok {
					t.Fatalf("FAIL: Proof verification failed for leaf '%s' at query %d", leaf, i+1)
				}
				okCount++
			}

			// Failed queries
			numFail := tc.numQueries - numSuccess
			for i := 0; i < numFail; i++ {
				leaf := "not_exist_leaf_" + strconv.Itoa(i)

				startProof := time.Now()
				proof, left, right := tree.GetProofByValue(leaf)
				proofTime := time.Since(startProof)

				startVerify := time.Now()
				ok := VerifyProofMMR(leaf, tree.GetRoot(leaf), proof, left, right)
				verifyTime := time.Since(startVerify)

				totalProof += proofTime
				totalVerify += verifyTime

				if ok {
					t.Fatalf("FAIL: Non-existent leaf '%s' should not verify at query %d", leaf, i+1+numSuccess)
				}
				okCount++
			}

			t.Logf("%d queries: total proof time %v (average %.3f µs), total verification time %v (average %.3f µs), correct: %d/%d",
				tc.numQueries, totalProof, float64(totalProof.Microseconds())/float64(tc.numQueries),
				totalVerify, float64(totalVerify.Microseconds())/float64(tc.numQueries),
				okCount, tc.numQueries)
		})
	}
}
