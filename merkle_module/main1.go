package main

import (
	"context"
	"database/sql"
	"log"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	"merkle_module/infra/storage"
	"merkle_module/merkletree"
	"merkle_module/utils"
	"os"
	"sync/atomic"
	"time"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"

	"github.com/ethereum/go-ethereum/common/lru"
)

func main() {
	databaseURL := getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/merkle_tree?sslmode=disable")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// prevent error: too many request...
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	db.SetConnMaxLifetime(time.Minute * 5)
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database")

	repo := storage.NewMerklePostgres(db)
	cache := lru.NewCache[string, *merkletree.MerkleTree](10) // issuerDID -> *MerkleTree
	// stress test go routine
	ctx := context.Background()
	issuerDID := "did:example:test_cli"
	numLeaves := 1000 // Change this to increase/decrease the number of leaves
	channel := make(chan *entities.MerkleNode, numLeaves)
	datas := make([][]byte, numLeaves)

	var count atomic.Int32
	for i := 0; i < numLeaves; i++ {
		data := getRandomData()
		datas[i] = data
		dataHash := utils.Hash(data)
		go func(data []byte) {
			node, err := addLeaf(ctx, repo, cache, issuerDID, data)
			if err != nil {
				log.Printf("Error adding leaf: %v", err)
				return
			}
			count.Add(1)
			channel <- node
		}(dataHash)
	}
	// verify the nodes added
	for i := 0; i < numLeaves; i++ {
		node := <-channel
		if node == nil {
			log.Println("Received nil node from channel")
			return
		}
		// Check if tree is cached
		tree, exists := cache.Get(issuerDID)
		if exists && tree.GetTreeID() == node.TreeID {
			if !verifyProof(tree, node, datas[i], i) {
				return
			}
			continue
		}
		nodes, err := repo.GetNodesByTreeID(ctx, node.TreeID)
		if err != nil {
			log.Printf("Error getting nodes by tree ID: %v", err)
			return
		}
		tree, err = merkletree.NewMerkleTree(nodes, node.TreeID)
		if err != nil {
			log.Printf("Error creating Merkle tree: %v", err)
			return
		}
		if !verifyProof(tree, node, datas[i], i) {
			return
		}
	}
}

func addLeaf(ctx context.Context, repo repo.Merkle, cache *lru.Cache[string, *merkletree.MerkleTree], issuerDID string, data []byte) (*entities.MerkleNode, error) {
	// check if the tree is in the cache
	tree, exists := cache.Get(issuerDID)
	if !exists {
		// if not, get the active tree for inserting
		nodes, err := repo.GetActiveTreeForInserting(ctx, issuerDID)
		if err != nil {
			return nil, err
		}
		// create a new Merkle tree
		tree, err = merkletree.NewMerkleTree(nodes.Nodes, nodes.TreeID)
		if err != nil {
			return nil, err
		}
		// add the tree to the cache
		cache.Add(issuerDID, tree)
	}

	tree.AddLeaf(data)

	// add to the database
	nodeID, err := tree.GetLastNodeID()
	_, err = repo.AddNode(ctx, tree.GetTreeID(), nodeID, data)
	if err != nil {
		return nil, err
	}

	return &entities.MerkleNode{
		NodeID: nodeID,
		TreeID: tree.GetTreeID(),
		Data:   data,
	}, nil
}

func getRandomData() []byte {
	// Generate random data for testing
	data := make([]byte, 32) // Example: 32 bytes of random data
	for i := range data {
		data[i] = byte(i % 256) // Simple pattern for demonstration
	}
	return data
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

func verifyProof(tree *merkletree.MerkleTree, node *entities.MerkleNode, data []byte, index int) bool {
	proof, err := tree.GetProof(node.NodeID)
	if err != nil {
		log.Printf("Error while getting proof for node %d: %v", node.NodeID, err)
		return false
	}
	if !utils.Verify(proof, tree.GetMerkleRoot(), data) {
		log.Printf("Proof verification failed for node %d: NodeID=%d, TreeID=%d", index, node.NodeID, node.TreeID)
		return false
	}
	log.Printf("Proof verified for node %d: NodeID=%d, TreeID=%d", index, node.NodeID, node.TreeID)
	return true
}
