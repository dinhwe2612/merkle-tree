package main

import (
	"context"
	"crypto/rand"
	"log"
	credential "merkle_module/smartcontract"
	"merkle_module/utils"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	_ "github.com/joho/godotenv/autoload" // Automatically load .env file
	mt "github.com/txaty/go-merkletree"
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

	// 1. Sinh dữ liệu ngẫu nhiên
	var leaves [][]byte
	for i := 0; i < 10; i++ {
		leaves = append(leaves, randomBytes(32))
	}
	log.Printf("Generated %d random leaves", len(leaves))

	// 2. Tạo Merkle tree bằng go-merkletree
	tree, err := mt.New(utils.GetTreeConfig(), utils.ToBlockDataFromByteArray(leaves))
	if err != nil {
		log.Fatalf("Failed to create Merkle tree: %v", err)
	}
	root := tree.Root
	log.Printf("Merkle root: %x", root)

	// 3. Lấy proof cho một leaf bất kỳ
	leafIndex := 3 // ví dụ lấy leaf thứ 4
	proof := tree.Proofs[leafIndex]
	if err != nil {
		log.Fatalf("Failed to get proof: %v", err)
	}
	log.Printf("Proof for leaf %d has %d siblings", leafIndex, len(proof.Siblings))

	// 4. Chuẩn hóa dữ liệu cho smart contract
	var root32 [32]byte
	copy(root32[:], root)
	var proof32 [][32]byte
	for _, sib := range proof.Siblings {
		var h [32]byte
		copy(h[:], sib)
		proof32 = append(proof32, h)
	}
	var leaf32 [32]byte
	copy(leaf32[:], leaves[leafIndex])

	issuers := []common.Address{common.HexToAddress("0x5A9cC578d5CC3Af41caf52c3818E77Af8Ff578D2")}
	treeIDs := []int{1}
	roots := [][32]byte{root32}

	// 5. Gọi BatchUpdate trên smart contract
	log.Println("Calling BatchUpdate...")
	err = smartContract.SendRoot(issuers, treeIDs, roots)
	if err != nil {
		log.Fatalf("BatchUpdate failed: %v", err)
	}
	log.Println("BatchUpdate success!")

	// 6. Verify trên smart contract
	log.Println("Verifying on smart contract...")
	isVerified, err := smartContract.Verify(
		context.Background(),
		issuers[0],
		treeIDs[0],
		leaf32,
		proof32,
	)
	if err != nil {
		log.Fatalf("Smart contract verify failed: %v", err)
	}
	log.Printf("Smart contract verify result: %v", isVerified)
	if isVerified {
		log.Println("✅ Smart contract verification successful!")
	} else {
		log.Println("❌ Smart contract verification failed!")
	}
}

func randomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("Failed to generate random bytes: %v", err)
	}
	return b
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
