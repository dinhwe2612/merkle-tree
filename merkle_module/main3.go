package main

import (
	"context"
	"database/sql"
	"log"
	"merkle_module/app/services"
	"merkle_module/infra/storage"
	"merkle_module/merkletree"
	"os"

	"github.com/ethereum/go-ethereum/common/lru"
	_ "github.com/joho/godotenv/autoload" // Load environment variables from .env file
	_ "github.com/lib/pq"                 // PostgreSQL driver
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

	root, err := merkleService.GetRoot(ctx, 406)
	if err != nil {
		log.Fatalf("Failed to get root: %v", err)
	}
	log.Printf("Merkle root: %x", root)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
