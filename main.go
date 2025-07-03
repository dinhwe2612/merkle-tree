package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/crypto"
	mt "github.com/txaty/go-merkletree"
)

// first define a data structure with Serialize method to be used as data block
type testData struct {
	data []byte
}

func (t *testData) Serialize() ([]byte, error) {
	return t.data, nil
}

// generate dummy data blocks
func generateRandBlocks() (blocks []mt.DataBlock) {
	hexStrings := []string{
		"1a8cd71aeb2aa2af4b47bc876cbc93f6bbee71af2f83ffc9a3c4e2c860e6eff0",
		"57580b9c14e7ed6fe1f1e6245fda0c26d7213a2904a1c32d55bf354eaac5ac39",
	}
	for i, hexStr := range hexStrings {
		data := make([]byte, 32)
		_, err := fmt.Sscanf(hexStr, "%x", &data)
		handleError(err)

		block := &testData{data: data}
		blocks = append(blocks, block)
		fmt.Printf("Block %d data: %x\n", i, block.data)
	}
	return
}

func main() {
	blocks := generateRandBlocks()
	// config with Keccak256Hash as the hash function
	config := &mt.Config{
		HashFunc: func(block []byte) ([]byte, error) {
			return crypto.Keccak256(block), nil
		},
		SortSiblingPairs: true,
	}
	// create a new Merkle Tree with the config and data blocks
	tree, err := mt.New(config, blocks)
	handleError(err)
	// get proofs
	proofs := tree.Proofs
	// or you can also verify the proofs without the tree but with Merkle root
	// obtain the Merkle root
	verified, err := mt.Verify(blocks[0], proofs[0], tree.Root, config)
	handleError(err)
	if !verified {
		fmt.Println("Proof verification failed for block 0")
		return
	}
	rootHash := tree.Root
	fmt.Printf("Merkle root: %x\n", rootHash)
	// print the block[0] hashes
	fmt.Printf("Block 0 hash: %x\n", tree.Leaves[0])
	// print proof[0]'s sibling hashes
	fmt.Println("Proof for leaf 0:")
	for _, sibling := range proofs[0].Siblings {
		fmt.Printf("%x\n", sibling)
	}
}

func hexStringToBytes(hexStr string) []byte {
	var data []byte
	_, err := fmt.Sscanf(hexStr, "%x", &data)
	if err != nil {
		fmt.Printf("Error converting hex string to bytes: %v\n", err)
		return nil
	}
	return data
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
