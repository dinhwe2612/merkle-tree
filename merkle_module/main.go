package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"merkle_module/app/interfaces"
	"merkle_module/app/services"
	"merkle_module/infra/storage"
	"os"

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

	// Create repository and service
	merkleRepo := storage.NewMerklePostgres(db)
	merkleService := services.NewMerkleService(merkleRepo)

	ctx := context.Background()
	testWithRandomData(ctx, merkleService, "did:example:123456789")
	testWithRandomData(ctx, merkleService, "did:example:987654321")
}

// generateRandomHash generates a random hash string
func generateRandomHash() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// testWithRandomData tests the merkle tree with random data
func testWithRandomData(ctx context.Context, merkleService interfaces.Merkle, issuerDID string) {
	fmt.Printf("=== Random Data Test for %s ===\n", issuerDID)
	// Generate random test data
	randomHashes := make([]string, 5)
	for i := 0; i < 5; i++ {
		randomHashes[i] = generateRandomHash()
	}

	fmt.Println("\n1. Adding random nodes to database...")
	for i, hash := range randomHashes {
		// Use service layer to add leaf
		err := merkleService.AddLeaf(ctx, issuerDID, []byte(hash))
		if err != nil {
			log.Printf("Failed to add random node %d: %v", i, err)
		} else {
			fmt.Printf("Added random node %d: %s\n", i, hash)
		}
	}

	// Test with a specific random hash
	testRandomHash := randomHashes[2]
	fmt.Printf("\n2. Testing with random hash: %s\n", testRandomHash)

	// Use service layer to get proof
	proof, err := merkleService.GetProof(ctx, issuerDID, []byte(testRandomHash))
	if err != nil {
		log.Printf("Failed to get proof for random hash: %v", err)
		return
	}
	fmt.Printf("Successfully generated proof for random hash\n")

	// Test proof verification
	if len(proof) > 0 {
		isValid, err := merkleService.VerifyProof(ctx, issuerDID, []byte(testRandomHash), proof)
		if err != nil {
			log.Printf("Failed to verify proof: %v", err)
		} else if isValid {
			fmt.Printf("Random hash proof verification successful!\n")
		} else {
			fmt.Printf("Random hash proof verification failed!\n")
		}

		// Test with a different random hash (should fail)
		differentHash := generateRandomHash()
		proof2, err := merkleService.GetProof(ctx, issuerDID, []byte(differentHash))
		if err == nil && len(proof2) > 0 {
			isValid2, err := merkleService.VerifyProof(ctx, issuerDID, []byte(differentHash), proof2)
			if err == nil && !isValid2 {
				fmt.Printf("Correctly rejected proof for different hash\n")
			} else {
				fmt.Printf("Should have rejected proof for different hash\n")
			}
		}
	}

	fmt.Println("\n=== Random Data Test completed ===\n")
}

// getEnv gets environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
