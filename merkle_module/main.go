package main

import (
	"context"
	"database/sql"
	"log"
	"merkle_module/app/services"
	"merkle_module/domain/entities"
	"merkle_module/infra/storage"
	"merkle_module/utils"
	"os"
	"time"

	"merkle_module/merkletree"

	"github.com/ethereum/go-ethereum/common/lru"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	connStr := getEnv("DATABASE_URL", "postgres://user:password@localhost:port/merkle_tree?sslmode=disable")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	merkleRepo := storage.NewMerklePostgres(db)
	merkleCache := lru.NewCache[int, *merkletree.MerkleTree](10)
	issuerCache := lru.NewCache[string, int](10)
	merkleService := services.NewMerkleService(merkleRepo, merkleCache, issuerCache)

	ctx := context.Background()
	issuerDID := "did:example:test_cli"
	numLeaves := 100
	channel := make(chan result, numLeaves)
	datas := make([][]byte, numLeaves)

	start := time.Now()
	failCount := 0
	addFailCount := 0
	for i := 0; i < numLeaves; i++ {
		data := randomBytes(32)
		datas[i] = data
		dataHash := utils.Hash(data)
		go func(idx int, data []byte) {
			node, err := merkleService.AddLeaf(ctx, issuerDID, data)
			if err != nil {
				log.Printf("Error adding leaf: %v", err)
				channel <- result{idx, nil}
				return
			}
			channel <- result{idx, node}
		}(i, dataHash)
	}

	// check for duplicate NodeID in the same TreeID
	type nodeKey struct {
		TreeID int
		NodeID int
	}
	nodeMap := make(map[nodeKey]int) // value: index of first occurrence
	for i := 0; i < numLeaves; i++ {
		res := <-channel
		if res.node == nil {
			log.Println("Received nil node from channel")
			addFailCount++
			failCount++
			continue
		}
		// check value in range [1, numLeaves]
		if res.node.NodeID < 1 || res.node.NodeID > numLeaves {
			log.Printf("Invalid NodeID: %d for index %d", res.node.NodeID, res.idx)
			addFailCount++
			failCount++
			continue
		}
		proof, err := merkleService.GetProof(ctx, res.node.TreeID, res.node.NodeID)
		if err != nil {
			log.Printf("Error getting proof for node %d: %v", res.node.NodeID, err)
			failCount++
			continue
		}
		root, err := merkleService.GetRoot(ctx, res.node.TreeID)
		if err != nil {
			log.Printf("Error getting root for node %d: %v", res.node.NodeID, err)
			failCount++
			continue
		}
		if !utils.Verify(proof, root, datas[res.idx]) {
			log.Printf("Proof verification failed for node %d: NodeID=%d, TreeID=%d", res.idx, res.node.NodeID, res.node.TreeID)
			failCount++
		} else {
			log.Printf("Proof verified for node %d: NodeID=%d, TreeID=%d", res.idx, res.node.NodeID, res.node.TreeID)
		}
		key := nodeKey{TreeID: res.node.TreeID, NodeID: res.node.NodeID}
		if idx, exists := nodeMap[key]; exists {
			log.Printf("Duplicate NodeID found: TreeID=%d, NodeID=%d (indexes: %d and %d)", res.node.TreeID, res.node.NodeID, idx, res.idx)
			return
		}
		nodeMap[key] = res.idx
	}
	elapsed := time.Since(start)
	log.Printf("Total add leaf fail: %d", addFailCount)
	log.Printf("Total verify proof fail: %d", failCount)
	log.Printf("Total elapsed time: %s", elapsed)

	// CASE: Add duplicate data
	log.Println("\n=== EDGE CASE 1: Add duplicate data ===")
	issuerDID = "did:example:test_cli_dup"
	dupData := randomBytes(32)
	for i := 0; i < 5; i++ {
		node, err := merkleService.AddLeaf(ctx, issuerDID, dupData)
		if err != nil {
			log.Printf("Add duplicate leaf error (iteration %d): %v", i, err)
		} else {
			log.Printf("Add duplicate leaf success (iteration %d): NodeID=%d, TreeID=%d", i, node.NodeID, node.TreeID)
		}
	}

	// CASE: Add MAX_LEAFS + 1 leaves to test tree full and new tree creation
	log.Println("\n=== EDGE CASE 2: Add MAX_LEAFS + 1 leaves to test tree full and new tree creation ===")
	issuerDID = "did:example:test_cli_max_leafs"
	maxLeafs := utils.MAX_LEAFS
	var lastTreeID, newTreeID int
	for i := 0; i < maxLeafs+1; i++ {
		data := randomBytes(32)
		node, err := merkleService.AddLeaf(ctx, issuerDID, data)
		if err != nil {
			log.Printf("Add leaf error at i=%d: %v", i, err)
			continue
		}
		if i == maxLeafs-1 {
			lastTreeID = node.TreeID
		}
		if i == maxLeafs {
			newTreeID = node.TreeID
		}
	}
	if lastTreeID != 0 && newTreeID != 0 {
		if lastTreeID != newTreeID {
			log.Printf("New tree created after reaching MAX_LEAFS: lastTreeID=%d, newTreeID=%d", lastTreeID, newTreeID)
		} else {
			log.Printf("TreeID did not change after MAX_LEAFS, possible bug: treeID=%d", lastTreeID)
		}
	} else {
		log.Printf("Could not determine treeID transition, lastTreeID=%d, newTreeID=%d", lastTreeID, newTreeID)
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 256)
	}
	return b
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

type result struct {
	idx  int
	node *entities.MerkleNode
}
