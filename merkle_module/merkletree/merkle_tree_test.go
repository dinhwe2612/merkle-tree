package merkletree

import (
	"fmt"
	"merkle_module/utils"
	"sync"
	"testing"
)

func TestMerkleTree(t *testing.T) {
	cid1 := "QmcPGZtdbUn9TT3rbDU6RkTkjx1DnJ3ENQUBhhoVF15KwA"
	cid2 := "QmVhveySRaD1fTSkZhdZ6MYWhFABgC4FhwtJme6dUJuW4u"
	cid3 := "QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco"
	cid4 := "QmbWqxBEKC3P8tqsKc98xmWNzrzDtRLMiMPL8wBuTGsMnR"

	data := [][]byte{[]byte(cid1), []byte(cid2), []byte(cid3)}

	hashData := make([][]byte, len(data))
	for i, cid := range data {
		hashData[i] = utils.Hash(cid)
		fmt.Printf("Hash of %s: %x\n", cid, hashData[i])
	}

	tree, err := NewMerkleTree(hashData, 0) // 0 is a placeholder for treeID, as we don't need it here
	if err != nil {
		t.Fatalf("Failed to create Merkle Tree: %v", err)
	}

	fmt.Printf("Merkle Tree created with %d leaves\n", tree.numLeafs)
	fmt.Printf("Root hash: %x\n", tree.nodes[1])

	// Check if the tree contains the original data
	for _, cid := range data {
		hashData := utils.Hash([]byte(cid))
		if !tree.Contains(hashData) {
			t.Errorf("Merkle Tree does not contain data: %s", cid)
			return
		} else {
			fmt.Printf("Merkle Tree contains data: %s\n", cid)
		}
	}

	// Verify the data
	for i, cid := range data {
		proof, err := tree.GetProof(i + 1) // Get proof for the leaf node, position starts from 1
		if err != nil {
			t.Errorf("Failed to get proof for data %s: %v", cid, err)
			continue
		}

		if !utils.Verify(proof, tree.nodes[1], []byte(cid)) {
			t.Errorf("Proof verification failed for data: %s", cid)
			return
		} else {
			fmt.Printf("Proof verified successfully for data: %s\n", cid)
		}
	}

	// Add a new node to the tree
	newData := []byte(cid4)
	newDataHash := utils.Hash(newData)
	tree.AddLeaf(newDataHash)
	fmt.Printf("Added new data: %s\n", cid4)
	if !tree.Contains(newDataHash) {
		t.Errorf("Merkle Tree does not contain newly added data: %s", cid4)
	} else {
		fmt.Printf("Merkle Tree contains newly added data: %s\n", cid4)
	}
	// Verify the new data
	proof, err := tree.GetProof(tree.numLeafs) // Get proof for the newly
	if err != nil {
		t.Errorf("Failed to get proof for new data %s: %v", cid4, err)
	} else {
		if !utils.Verify(proof, tree.nodes[1], newData) {
			t.Errorf("Proof verification failed for new data: %s", cid4)
		} else {
			fmt.Printf("Proof verified successfully for new data: %s\n", cid4)
		}
	}

}

func TestStress(t *testing.T) {
	// This test is designed to stress the Merkle Tree implementation
	maxLeafs := min(1<<15, utils.MAX_LEAFS) // 32768 leaves
	tree, err := NewMerkleTree(nil, maxLeafs)
	if err != nil {
		t.Fatalf("Failed to create Merkle Tree: %v", err)
	}

	fmt.Printf("Stress Test: Created Merkle Tree with max %d leaves\n", maxLeafs)
	for i := 0; i < maxLeafs; i++ {
		data := fmt.Sprintf("data-%d", i)
		hashData := utils.Hash([]byte(data))
		tree.AddLeaf(hashData)
		if !tree.Contains(hashData) {
			t.Errorf("Merkle Tree does not contain data: %s", data)
			return
		}
	}
	fmt.Printf("Stress Test: Added %d leaves to the Merkle Tree\n", maxLeafs)
	root := tree.GetMerkleRoot()
	if len(root) == 0 {
		t.Errorf("Merkle Tree root is empty after adding %d leaves", maxLeafs)
		return
	}
	fmt.Printf("Stress Test: Merkle Tree root hash after adding %d leaves: %x\n", maxLeafs, root)
	for i := 0; i < maxLeafs; i++ {
		data := fmt.Sprintf("data-%d", i)
		proof, err := tree.GetProof(i + 1) // Get proof for the leaf node, position starts from 1
		if err != nil {
			t.Errorf("Failed to get proof for data %s: %v", data, err)
			continue
		}

		if !utils.Verify(proof, root, []byte(data)) {
			t.Errorf("Proof verification failed for data: %s", data)
			return
		}
	}
	fmt.Printf("Stress Test: All proofs verified successfully for %d leaves\n", maxLeafs)
}

func TestConcurrent(t *testing.T) {
	// This test is designed to test the Merkle Tree implementation with concurrent access
	tree, err := NewMerkleTree(nil, 0) // 0 is a placeholder for treeID, as we don't need it here
	if err != nil {
		t.Fatalf("Failed to create Merkle Tree: %v", err)
	}
	fmt.Printf("Concurrent Test: Created Merkle Tree with %d leaves\n", tree.maxLeafs)
	type result struct {
		idx    int
		NodeID int
		data   []byte
	}
	channel := make(chan result, tree.maxLeafs)
	var wg sync.WaitGroup
	wg.Add(tree.maxLeafs)
	for i := 0; i < tree.maxLeafs; i++ {
		data := fmt.Sprintf("data-%d", i)
		go func(idx int, data string) {
			hashData := utils.Hash([]byte(data))
			nodeID := tree.AddLeaf(hashData)
			if nodeID < 0 {
				channel <- result{idx, -1, nil}
				return
			}
			channel <- result{idx, nodeID, []byte(data)}
			wg.Done()
		}(i, data)
	}
	// wait for all goroutines to finish
	wg.Wait()

	for i := 0; i < tree.maxLeafs; i++ {
		res := <-channel
		if res.NodeID < 0 {
			t.Errorf("Failed to add leaf for index %d", res.idx)
			return
		}
		hashData := utils.Hash(res.data)
		if !tree.Contains(hashData) {
			t.Errorf("Merkle Tree does not contain data for index %d: %s", res.idx, res.data)
			return
		}
		proof, err := tree.GetProof(res.NodeID)
		if err != nil {
			t.Errorf("Failed to get proof for index %d: %v", res.idx, err)
			return
		}
		if !utils.Verify(proof, tree.GetMerkleRoot(), res.data) {
			t.Errorf("Proof verification failed for index %d: %s", res.idx, res.data)
			return
		}
	}
}
