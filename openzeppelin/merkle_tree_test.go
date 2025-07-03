package openzeppelin

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

func TestMerkleTree(t *testing.T) {
	cid1 := []byte("QmcPGZtdbUn9TT3rbDU6RkTkjx1DnJ3ENQUBhhoVF15KwA")
	cid2 := []byte("QmVhveySRaD1fTSkZhdZ6MYWhFABgC4FhwtJme6dUJuW4u")
	cid3 := []byte("QmXoypizjW3WknFiJnKLwHCnL72vedxjQkDDP1mXWo6uco")
	cid4 := []byte("QmbWqxBEKC3P8tqsKc98xmWNzrzDtRLMiMPL8wBuTGsMnR")

	data := [][]byte{cid1, cid2, cid3}

	tree, err := NewMerkleTree(data)
	if err != nil {
		t.Fatalf("Failed to create Merkle Tree: %v", err)
	}

	initialRoot := tree.GetMerkleRoot()
	fmt.Printf("Initial Merkle Root: %x\n", initialRoot)

	t.Log("--- Adding a new leaf ---")
	err = tree.AddLeaf(cid4)
	if err != nil {
		t.Fatalf("Failed to add leaf: %v", err)
	}

	updatedRoot := tree.GetMerkleRoot()
	fmt.Printf("Updated Merkle Root: %x\n", updatedRoot)

	leafToProve := cid1
	proof, err := tree.GetProof(leafToProve)
	if err != nil {
		t.Fatalf("Failed to get proof: %v", err)
	}

	fmt.Printf("\n--- Proof for leaf %s ---\n", leafToProve)
	fmt.Printf("Leaf to prove: 0x%x\n", crypto.Keccak256(leafToProve))
	isValid := Verify(proof, updatedRoot, leafToProve)
	if !isValid {
		t.Errorf("Proof is invalid for leaf %s", leafToProve)
	} else {
		fmt.Println("Proof is valid!")
	}
}

func TestStress(t *testing.T) {
	// generate bytes
	data := make([][]byte, 1000)
	for i := 0; i < 1000; i++ {
		data[i] = []byte(fmt.Sprintf("data-%d", i))
	}
	tree, err := NewMerkleTree(data)
	if err != nil {
		t.Fatalf("Failed to create Merkle Tree: %v", err)
	}
	root := tree.GetMerkleRoot()
	// add new leafs
	for i := 1000; i < 2000; i++ {
		data = append(data, []byte(fmt.Sprintf("data-%d", i)))
		err := tree.AddLeaf(data[i])
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
		proof, err := tree.GetProof(leaf)
		if err != nil {
			t.Fatalf("Failed to get proof for leaf %d: %v", leafIndex, err)
		}
		isValid := Verify(proof, newRoot, leaf)
		if !isValid {
			t.Errorf("Proof is invalid for leaf %s", leaf)
		} else {
			fmt.Printf("Proof is valid for leaf %s\n", leaf)
		}
	}
	// verify proof for a leaf that does not exist
	nonExistentLeaf := []byte("non-existent-leaf")
	_, err = tree.GetProof(nonExistentLeaf)
	if err == nil {
		t.Errorf("Expected error when getting proof for non-existent leaf, but got none")
	} else {
		fmt.Printf("Correctly received error for non-existent leaf: %v\n", err)
	}
}
