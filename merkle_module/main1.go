package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"log"
	"merkle_module/app/services"
	"merkle_module/infra/storage"
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

	// add data into the merkle tree 100 leaves every 2s
	ctx := context.Background()
	issuerDIDs := []string{
		"did:example:test_cli",
		"did:example:test_cli_2",
		"did:example:test_cli_3",
	}
	for {
		for i := 0; i < 10000; i++ {
			for _, issuerDID := range issuerDIDs {
				go func(issuerDID string) {
					data := randomBytes(32)
					// make a copy of the data
					dataCopy := make([]byte, len(data))
					copy(dataCopy, data)
					node, err := merkleService.AddLeaf(ctx, issuerDID, dataCopy)
					if err != nil {
						log.Printf("Error adding leaf: %v", err)
						return
					}
					if node == nil {
						log.Printf("Error: node is nil after adding leaf")
						return
					}
					log.Printf("Added leaf for issuer %s", issuerDID)
				}(issuerDID)
			}
		}
		log.Println("Waiting for 2 seconds before adding more leaves...")
		time.Sleep(2 * time.Second)
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		if _, err := rand.Read(b[i : i+1]); err != nil {
			log.Fatalf("Failed to generate random bytes: %v", err)
		}
	}
	return b
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
