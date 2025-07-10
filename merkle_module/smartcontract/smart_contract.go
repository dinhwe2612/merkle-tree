package credential

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SmartContract struct {
	client          *ethclient.Client
	contract        *Credential
	auth            *bind.TransactOpts
	contractAddress common.Address
}

func NewEthClient(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	return client, nil
}

func NewTransactOpts(client *ethclient.Client, privateKey string, chainID *big.Int) (*bind.TransactOpts, error) {
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor: %w", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to suggest gas price: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // no ether being sent
	auth.GasLimit = uint64(3000000) // adjust as needed
	auth.GasPrice = gasPrice

	return auth, nil
}

func NewSmartContract(client *ethclient.Client, contract *Credential, contractAddress common.Address, auth *bind.TransactOpts) *SmartContract {
	return &SmartContract{
		client:          client,
		contract:        contract,
		auth:            auth,
		contractAddress: contractAddress,
	}
}

func (sc *SmartContract) SendRoot(issuers []common.Address, treeIDs []int, roots [][32]byte) error {
	// int to big.Int conversion for treeIDs
	treeIDsBig := make([]*big.Int, len(treeIDs))
	for i, id := range treeIDs {
		treeIDsBig[i] = big.NewInt(int64(id))
	}

	// Send the Merkle root to the smart contract
	tx, err := sc.contract.BatchUpdateTreeRoots(
		sc.auth,
		issuers,
		treeIDsBig,
		roots,
	)
	if err != nil {
		return fmt.Errorf("failed to send root to contract: %w", err)
	}
	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(context.Background(), sc.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	fmt.Printf("Transaction mined! Block number: %d\n", receipt.BlockNumber)

	return nil
}

func (sc *SmartContract) Verify(ctx context.Context, from common.Address, issuer common.Address, treeIndex int, leaf [32]byte, proof [][32]byte) (bool, error) {
	abiBytes, err := os.ReadFile("./smartcontract/build/NDACredential.abi")
	if err != nil {
		log.Fatalf("Failed to read ABI file: %v", err)
	}
	contractABI, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	var proofHash []common.Hash
	for _, p := range proof {
		proofHash = append(proofHash, common.BytesToHash(p[:]))
	}

	callData, err := contractABI.Pack("verifyVC", issuer, big.NewInt(int64(treeIndex)), leaf, proofHash)
	if err != nil {
		log.Fatalf("Failed to pack data for contract call: %v", err)
	}

	result, err := sc.client.CallContract(context.Background(), ethereum.CallMsg{
		From: sc.auth.From,
		To:   &sc.contractAddress,
		Data: callData,
	}, nil)
	if err != nil {
		log.Fatalf("Failed to call contract: %v", err)
	}

	var isVerified bool
	err = contractABI.UnpackIntoInterface(&isVerified, "verifyVC", result)
	if err != nil {
		log.Fatalf("Failed to unpack contract result: %v", err)
	}

	return isVerified, nil
}
