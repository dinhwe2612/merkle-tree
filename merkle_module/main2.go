package main

import (
	"context"
	"database/sql"
	"log"
	"merkle_module/cronjob"
	"merkle_module/infra/storage"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload" // Load environment variables from .env file
	_ "github.com/lib/pq"                 // PostgreSQL driver
)

func main() {
	// Load environment variables
	connStr := getEnv("DATABASE_URL", "postgres://user:password@localhost:port/merkle_tree?sslmode=disable")
	// Connect to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// Initialize context
	ctx := context.Background()
	// Initialize Merkle repository
	merkleRepo := storage.NewMerklePostgres(db)
	// CASE: Use cron job to sync Merkle root
	log.Println("\n=== CASE 1: Use cron job to sync Merkle root ===")
	syncJob := cronjob.NewAsyncJob(ctx, merkleRepo)
	syncJob.Start()
	log.Println("Cron job started to sync Merkle root")
	for len(syncJob.GetRunningJobs()) > 0 {
		time.Sleep(100 * time.Second) // Wait for the job to complete
	}
	syncJob.Stop()
	log.Println("Cron job stopped")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
