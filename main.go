package main

import (
	"crypto/rand"
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
func generateRandBlocks(size int) (blocks []mt.DataBlock) {
	for i := 0; i < size; i++ {
		block := &testData{
			data: make([]byte, 32),
		}
		_, err := rand.Read(block.data)
		handleError(err)
		blocks = append(blocks, block)
	}
	return
}

func main() {
	blocks := generateRandBlocks(2)
	// config with Keccak256Hash as the hash function
	config := &mt.Config{
		HashFunc: func(block []byte) ([]byte, error) {
			if len(block) != 32 && len(block) != 64 {
				return nil, fmt.Errorf("invalid block size: %d, expected 32 or 64 bytes", len(block))
			}
			hash := crypto.Keccak256(block)
			// hash = crypto.Keccak256(hash)
			return hash, nil
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
	fmt.Println()
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
