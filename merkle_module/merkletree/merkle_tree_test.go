package merkletree

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestMerkleTree(t *testing.T) {
	cid1 := "QmcPGZtdbUn9TT3rbDU6RkTkjx1DnJ3ENQUBhhoVF15KwA"
	cid2 := "QmVhveySRaD1fTSkZhdZ6MYWhFABgC4FhwtJme6dUJuW4u"
	cid3 := "QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco"
	cid4 := "QmbWqxBEKC3P8tqsKc98xmWNzrzDtRLMiMPL8wBuTGsMnR"

	data := []string{cid1, cid2, cid3}

	tree, err := NewMerkleTree(data)
	if err != nil {
		t.Fatalf("Failed to create Merkle Tree: %v", err)
	}

	initialRoot := tree.GetMerkleRoot()
	fmt.Printf("Initial Merkle Root: %x\n", initialRoot)

	t.Log("--- Adding a new leaf ---")
	err = tree.AddLeaf([]byte(cid4))
	if err != nil {
		t.Fatalf("Failed to add leaf: %v", err)
	}

	updatedRoot := tree.GetMerkleRoot()
	fmt.Printf("Updated Merkle Root: %x\n", updatedRoot)

	leafToProve := cid1
	proof, err := tree.GetProof([]byte(leafToProve))
	if err != nil {
		t.Fatalf("Failed to get proof: %v", err)
	}

	fmt.Printf("\n--- Proof for leaf %s ---\n", leafToProve)
	fmt.Printf("Leaf to prove: 0x%x\n", crypto.Keccak256([]byte(leafToProve)))
	isValid := Verify(proof, updatedRoot, []byte(leafToProve))
	if !isValid {
		t.Errorf("Proof is invalid for leaf %s", leafToProve)
	} else {
		fmt.Println("Proof is valid!")
	}
}

func TestStress(t *testing.T) {
	// generate strings
	data := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = fmt.Sprintf("data-%d", i)
	}
	tree, err := NewMerkleTree(data)
	if err != nil {
		t.Fatalf("Failed to create Merkle Tree: %v", err)
	}
	root := tree.GetMerkleRoot()
	// add new leafs
	for i := 1000; i < 2000; i++ {
		newData := fmt.Sprintf("data-%d", i)
		data = append(data, newData)
		err := tree.AddLeaf([]byte(newData))
		if err != nil {
			t.Fatalf("Failed to add leaf: %v", err)
		}
	}
	newRoot := tree.GetMerkleRoot()
	if bytes.Compare(root, newRoot) == 0 {
		t.Errorf("Root should have changed after adding new leaves")
	} else {
		fmt.Printf("New Merkle Root after adding 1000 leaves: %x\n", newRoot)
	}
	// verify proof by randomly selecting a leaf
	for i := 0; i < 2000; i++ {
		leafIndex := i % len(data)
		leaf := data[leafIndex]
		proof, err := tree.GetProof([]byte(leaf))
		if err != nil {
			t.Fatalf("Failed to get proof for leaf %d: %v", leafIndex, err)
		}
		isValid := Verify(proof, newRoot, []byte(leaf))
		if !isValid {
			t.Errorf("Proof is invalid for leaf %s", leaf)
			return
		}
	}
	// verify proof for a leaf that does not exist
	nonExistentLeaf := []byte("non-existent-leaf")
	_, err = tree.GetProof(nonExistentLeaf)
	if err == nil {
		t.Errorf("Expected error when getting proof for non-existent leaf, but got none")
	}
}
