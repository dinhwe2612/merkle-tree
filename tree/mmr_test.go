package tree

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestMMR(t *testing.T) {
	for _, tt := range TestCases {
		t.Run(tt.name, func(t *testing.T) {
			tree := &MMR{}

			leaves := make([]string, tt.numLeaves)
			for i := 0; i < tt.numLeaves; i++ {
				leaves[i] = "leaf_" + strconv.Itoa(i)
			}

			startAdd := time.Now()
			for _, v := range leaves {
				tree.AddLeaf(v)
			}
			addTime := time.Since(startAdd)
			t.Logf("Added %d leaves in: %v (avg %.3f µs/leaf)", tt.numLeaves, addTime, float64(addTime.Microseconds())/float64(tt.numLeaves))

			rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
			totalProof := time.Duration(0)
			totalVerify := time.Duration(0)
			okCount := 0

			numSuccessful := int(float64(tt.numQueries) * 0.8)
			numFailed := tt.numQueries - numSuccessful

			// Successful queries
			for i := 0; i < numSuccessful; i++ {
				leaf := leaves[rnd.Intn(tt.numLeaves)]
				startProof := time.Now()
				peakPath, leftPeaks, rightPeaks := tree.GetProof(leaf)
				proofTime := time.Since(startProof)

				startVerify := time.Now()
				ok := VerifyProofMMR(leaf, tree.GetRoot(), peakPath, leftPeaks, rightPeaks)
				verifyTime := time.Since(startVerify)

				if !ok {
					t.Fatalf("Proof verification failed for existing leaf '%s'", leaf)
				}
				totalProof += proofTime
				totalVerify += verifyTime
				okCount++
			}

			// Failed queries
			for i := 0; i < numFailed; i++ {
				leaf := "leaf_" + strconv.Itoa(tt.numLeaves+rnd.Int())
				startProof := time.Now()
				peakPath, _, _ := tree.GetProof(leaf)
				proofTime := time.Since(startProof)

				if peakPath != nil && len(peakPath) > 0 {
					t.Fatalf("Expected no proof for non-existent leaf '%s'", leaf)
				}
				totalProof += proofTime
				okCount++
			}

			t.Logf("%d queries: total proof time %v (avg %.3f µs), total verification time %v (avg %.3f µs), correct: %d/%d",
				tt.numQueries,
				totalProof, float64(totalProof.Microseconds())/float64(tt.numQueries),
				totalVerify, float64(totalVerify.Microseconds())/float64(tt.numQueries),
				okCount, tt.numQueries)
		})
	}
}
