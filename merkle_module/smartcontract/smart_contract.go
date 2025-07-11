package credential

import (
	"context"
	"fmt"
	"math/big"

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
	// Convert to big.Int
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
	if sc.contract == nil {
		return false, fmt.Errorf("smart contract is not initialized")
	}

	// Convert treeIndex to big.Int
	treeIndexBig := big.NewInt(int64(treeIndex))

	// Call the Verify function on the smart contract
	isVerified, err := sc.contract.VerifyVC(
		&bind.CallOpts{Context: ctx},
		issuer,
		treeIndexBig,
		leaf,
		proof,
	)
	if err != nil {
		return false, fmt.Errorf("failed to verify on smart contract: %w", err)
	}

	return isVerified, nil
}
