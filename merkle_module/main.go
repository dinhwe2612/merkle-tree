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
	numLeaves := 1000
	channel := make(chan *entities.MerkleNode, numLeaves)
	datas := make([][]byte, numLeaves)

	for i := 0; i < numLeaves; i++ {
		data := randomBytes(32)
		datas[i] = data
		dataHash := utils.Hash(data)
		go func(data []byte) {
			node, err := merkleService.AddLeaf(ctx, issuerDID, data)
			if err != nil {
				log.Printf("Error adding leaf: %v", err)
				channel <- nil
				return
			}
			channel <- node
		}(dataHash)
	}

	for i := 0; i < numLeaves; i++ {
		node := <-channel
		if node == nil {
			log.Println("Received nil node from channel")
			continue
		}
		proof, err := merkleService.GetProof(ctx, node.TreeID, node.NodeID)
		if err != nil {
			log.Printf("Error getting proof for node %d: %v", node.NodeID, err)
			continue
		}
		root, err := merkleService.GetRoot(ctx, node.TreeID)
		if err != nil {
			log.Printf("Error getting root for node %d: %v", node.NodeID, err)
			continue
		}
		if !utils.Verify(proof, root, datas[i]) {
			log.Printf("Proof verification failed for node %d: NodeID=%d, TreeID=%d", i, node.NodeID, node.TreeID)
		} else {
			log.Printf("Proof verified for node %d: NodeID=%d, TreeID=%d", i, node.NodeID, node.TreeID)
		}
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
