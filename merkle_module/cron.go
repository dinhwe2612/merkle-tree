package main

import (
	"context"
	"database/sql"
	"log"
	"merkle_module/cronjob"
	"merkle_module/infra/storage"
	credential "merkle_module/smartcontract"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
	ethClient, err := ethclient.Dial(getEnv("ETHEREUM_URL", ""))
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	contractAddress := common.HexToAddress(getEnv("CONTRACT_ADDRESS", ""))
	contract, err := credential.NewCredential(contractAddress, ethClient)
	if err != nil {
		log.Fatalf("Failed to instantiate contract: %v", err)
	}

	privateKey := getEnv("ACCOUNT_PRIVATE_KEY", "")
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	smartContract := credential.NewSmartContract(ethClient, contract, contractAddress, key, context.Background())
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
	syncJob := cronjob.NewAsyncJob(ctx, merkleRepo, smartContract)
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
