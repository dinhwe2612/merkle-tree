package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"log"
	"merkle_module/app/services"
	"merkle_module/infra/storage"
	"os"
	"sync"
	"time"

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
	merkleService := services.NewMerkleService(merkleRepo)

	// add data into the merkle tree 100 leaves every 2s
	ctx := context.Background()
	issuerDIDs := []string{
		"did:example:test_cli_1",
		"did:example:test_cli_2",
		"did:example:test_cli_3",
	}
	wait := sync.WaitGroup{}
	log.Println("Starting to add leaves to the Merkle tree...")
	MAX_ADD := 100
	wait.Add(len(issuerDIDs) * MAX_ADD)
	for {
		for i := 0; i < MAX_ADD; i++ {
			for _, issuerDID := range issuerDIDs {
				go func(issuerDID string) {
					data := randomBytes(32)
					// make a copy of the data
					dataCopy := make([]byte, len(data))
					copy(dataCopy, data)
					err := merkleService.AddLeaf(ctx, issuerDID, dataCopy)
					if err != nil {
						log.Printf("Error adding leaf: %v", err)
						return
					}
					log.Printf("Added leaf for issuer %s", issuerDID)
					wait.Done()
				}(issuerDID)
			}
		}
		log.Println("Waiting for 10 seconds before adding more leaves...")
		time.Sleep(10 * time.Second)
		break
	}
	wait.Wait()
	log.Println("All leaves added to the Merkle tree.")
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
