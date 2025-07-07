package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"merkle_module/app/interfaces"
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
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	merkleRepo := storage.NewMerklePostgres(db)
	merkleCache := lru.NewCache[int, *merkletree.MerkleTree](10)
	merkleService := services.NewMerkleService(merkleRepo, merkleCache)

	ctx := context.Background()
	issuerDID := "did:example:test_cli"
	data := []byte("test_data_1")

	fmt.Println("\n==== TEST AddLeaf ====")
	node, err := merkleService.AddLeaf(ctx, issuerDID, data)
	if err != nil {
		fmt.Printf("AddLeaf error: %v\n", err)
		return
	}
	fmt.Printf("AddLeaf success: NodeID=%d, TreeID=%d\n", node.NodeID, node.TreeID)

	fmt.Println("\n==== TEST GetProof ====")
	proof, err := merkleService.GetProof(ctx, node.TreeID, node.NodeID)
	if err != nil {
		fmt.Printf("GetProof error: %v\n", err)
		return
	}
	fmt.Printf("GetProof success: proof length=%d\n", len(proof))

	fmt.Println("\n==== TEST GetRoot ====")
	root, err := merkleService.GetRoot(ctx, node.TreeID)
	if err != nil {
		fmt.Printf("GetRoot error: %v\n", err)
		return
	}
	fmt.Printf("GetRoot success: %x\n", root)

	fmt.Println("\n==== TEST lỗi GetProof với nodeID không tồn tại ====")
	_, err = merkleService.GetProof(ctx, node.TreeID, 99999)
	if err != nil {
		fmt.Printf("GetProof (nodeID không tồn tại) error: %v\n", err)
	} else {
		fmt.Println("GetProof (nodeID không tồn tại) không báo lỗi (KHÔNG ĐÚNG)")
	}

	fmt.Println("\n==== STRESS TEST ====")
	numLeaves := 10
	stressTres(ctx, merkleService, issuerDID, numLeaves)
	fmt.Println("Stress test completed.")
}

func stressTres(ctx context.Context, service interfaces.Merkle, issuerDID string, numLeaves int) {
	var datas [][]byte
	var nodes []*entities.MerkleNode
	log.Printf("Starting stress test with %d leaves for issuer DID: %s", numLeaves, issuerDID)
	for i := 0; i < numLeaves; i++ {
		data := randomBytes(32) // Random 32 bytes
		node, err := service.AddLeaf(ctx, issuerDID, utils.Hash(data))
		if err != nil {
			log.Printf("AddLeaf error: %v", err)
			return
		}
		nodes = append(nodes, node)
		datas = append(datas, data)
	}
	fmt.Printf("Total leaves added: %d\n", len(datas))

	for i, node := range nodes {
		proof, err := service.GetProof(ctx, node.TreeID, node.NodeID)
		if err != nil {
			log.Printf("GetProof error for node %d: %v", i, err)
			return
		}
		root, err := service.GetRoot(ctx, node.TreeID)
		if err != nil {
			log.Printf("GetRoot error for node %d: %v", i, err)
			return
		}
		// Verify the proof
		if !utils.Verify(proof, root, datas[i]) {
			fmt.Printf("Proof verification failed for node %d: NodeID=%d, TreeID=%d\n", i, node.NodeID, node.TreeID)
			return
		}
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 256) // Just a simple pattern for testing
	}
	return b
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
