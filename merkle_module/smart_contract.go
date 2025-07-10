package main

import (
	"context"
	"encoding/hex"
	"log"
	credential "merkle_module/smartcontract"
	"merkle_module/utils"
	"os"

	"github.com/ethereum/go-ethereum/common"
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

	chainID, err := ethClient.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	auth, err := credential.NewTransactOpts(ethClient, getEnv("ACCOUNT_PRIVATE_KEY", ""), chainID)
	if err != nil {
		log.Fatalf("Failed to create transaction options: %v", err)
	}

	smartContract := credential.NewSmartContract(ethClient, contract, contractAddress, auth)

	issuers := []common.Address{
		common.HexToAddress("0x5A9cC578d5CC3Af41caf52c3818E77Af8Ff578D2"),
	}
	treeIDs := []int{1}
	// Merkle root as []byte
	rootHex := "9b98d7dd688546223e3e4f564102aceeb58d0607699d3da8bfc59014f4c18445"
	root, _ := hex.DecodeString(rootHex)
	roots := [][32]byte{[32]byte(root)}

	// Merkle proof as [][32]byte
	proofHex := []string{
		"77bf017d3c7c57b13f4075b398f072fc239d6e184b9b72454b759d33070dac49",
		"3b4d86955e34e3c7a99d614089a837d7c2f7cd58bf7fcea6e3ef2b53f711a5af",
		"13af5ac29c808ca3a5e7e7f379003de438dfd633558607bd4d901601f5d07f80",
		"cadd59817a3782796d6f24aad001d6733861674c1965ae5414da01f603ed9970",
		"c93dc07e2a20fabd2cfd1314849c9229ca4239ec5b348d88aeedd64cee158c25",
		"18b225d0260566572c7f754f6459979d423a4ae3276bffcbcb97cb420a163e98",
		"a8a885e3565ef4a865deb62414cb9fa563bb16919541a9aff64a930061bc316b",
		"8d76dc22a7f88007974d80737c5a541653e957df0b7bd3b55382c3e94a7ddcac",
		"c2d3ef9b061238bc0fd0ffb8a66f39be07a68c7baaf33a91b20cd8e2510ad589",
		"e9b1179138035c78862611fceddf63714c66e8bada31986ac4e78518e6c6d6c6",
		"17effb9cf3907d28704793e383f7b1da39fa690b1af1ca1166daf2c1d6a84bea",
		"9941170cb992c9e70fe3713f53565b275f4e53bc6c07236ac6fe8300760dbd63",
		"6413a91d50beed8e421030b4b870e0ae80a0a2ec5881e02fc46a45148b49c987",
		"835047751b64ed46147a26df6b0daf2aafdd72ec553a1b1762e9989a97585629",
		"82122931a6eea02f4b6e4baa216bd1dcbccdcb32a91f38d401440ff9451c6e9b",
		"2356399bf3663561e3318a42ed94bfa534a378ce06a52fa717e9c62d377fdcfc",
		"24031756c18b2f37f930a99e43a84a43bcc6c3b0ededee6eeedba33f41306fdd",
		"82cb4388380079a748ab856d28fc01e0b84b42f44198fea956d1174046752e2f",
		"fe18e90110b664fa6a688e32234e32552bae0ff06d7a0752a71dc93f5eed218f",
		"8f12a156d36736a4fb4e08160b7bdf7cf9e9efeb942f77f3b7584ca755980caa",
	}

	var proof [][32]byte
	for _, hexStr := range proofHex {
		bytes, _ := hex.DecodeString(hexStr)
		var hash [32]byte
		copy(hash[:], bytes)
		proof = append(proof, hash)
	}

	log.Printf("issuers: %v", issuers)
	log.Printf("treeIDs: %v", treeIDs)
	for i, r := range roots {
		log.Printf("roots[%d]: %x", i, r)
	}

	err = smartContract.SendRoot(issuers, treeIDs, roots)
	if err != nil {
		log.Fatalf("SendRoot failed: %v", err)
	} else {
		log.Println("SendRoot success!")
	}

	leafData := []byte("QmVhveySRaD1fTSkZhdZ6MYWhFABgC4FhwtJme6dUJuW4u")
	leafHash := utils.Hash(leafData)
	leaf32 := [32]byte{}
	copy(leaf32[:], leafHash)

	isVerified, err := smartContract.Verify(
		context.Background(),
		issuers[0],
		issuers[0],
		treeIDs[0],
		leaf32,
		proof,
	)
	if err != nil {
		log.Fatalf("Failed to verify: %v", err)
	}
	log.Printf("Verification result: %v", isVerified)
	if isVerified {
		log.Println("The leaf is verified successfully!")
	} else {
		log.Println("The leaf verification failed.")
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
