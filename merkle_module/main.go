package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
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

	merkleRepo := storage.NewMerklePostgres(db)
	merkleCache := storage.NewMerklesInMemory()
	merkleService := services.NewMerkleService(merkleRepo, merkleCache)

	ctx := context.Background()

	issuerDID := "did:example:test_" + generateRandomHash()[:8]
	testData := []byte(generateRandomData())

	fmt.Printf("\nMerkle Service Test\n==========================\n")
	fmt.Printf("Issuer: %s\n", issuerDID)
	fmt.Printf("Data:   %s\n", string(testData))

	// Add leaf
	fmt.Print("Adding leaf... ")
	merkleNode, err := merkleService.AddLeaf(ctx, issuerDID, testData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Success (Node ID: %d, Tree ID: %d)\n", merkleNode.NodeID, merkleNode.TreeID)

	// Try to add the same data again (should fail)
	fmt.Print("Adding duplicate data... ")
	_, err = merkleService.AddLeaf(ctx, issuerDID, testData)
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	} else {
		fmt.Println("Unexpected success - should have failed")
		return
	}

	// Add different data for the same issuer
	fmt.Print("Adding different data for same issuer... ")
	differentData := []byte(generateRandomData())
	fmt.Printf("Data: %s\n", string(differentData))
	merkleNode2, err := merkleService.AddLeaf(ctx, issuerDID, differentData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Success (Node ID: %d, Tree ID: %d)\n", merkleNode2.NodeID, merkleNode2.TreeID)

	// Get proof for first data
	fmt.Print("Getting proof for first data... ")
	proof, err := merkleService.GetProof(ctx, issuerDID, testData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Success (proof length: %d)\n", len(proof))

	// Verify proof for first data
	fmt.Print("Verifying proof for first data... ")
	isValid, err := merkleService.VerifyProof(ctx, issuerDID, testData, proof)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	if isValid {
		fmt.Println("Proof is valid!")
	} else {
		fmt.Println("Proof is invalid!")
	}

	// Get proof for second data
	fmt.Print("Getting proof for second data... ")
	proof2, err := merkleService.GetProof(ctx, issuerDID, differentData)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("Success (proof length: %d)\n", len(proof2))

	// Verify proof for second data
	fmt.Print("Verifying proof for second data... ")
	isValid2, err := merkleService.VerifyProof(ctx, issuerDID, differentData, proof2)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	if isValid2 {
		fmt.Println("Proof is valid!")
	} else {
		fmt.Println("Proof is invalid!")
	}
}

func generateRandomHash() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func generateRandomData() string {
	return fmt.Sprintf("test_data_%s", generateRandomHash()[:8])
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
