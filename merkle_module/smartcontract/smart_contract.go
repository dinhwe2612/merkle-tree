package credential

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type SmartContract struct {
	client          *ethclient.Client
	contract        *Credential
	contractAddress common.Address
	privateKey      *ecdsa.PrivateKey
	ctx             context.Context
}

func NewEthClient(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}
	return client, nil
}

func NewSmartContract(client *ethclient.Client, contract *Credential, contractAddress common.Address, privateKey *ecdsa.PrivateKey, ctx context.Context) *SmartContract {
	return &SmartContract{
		ctx:             ctx,
		client:          client,
		contract:        contract,
		contractAddress: contractAddress,
		privateKey:      privateKey,
	}
}

func (sc *SmartContract) SendRoot(issuers []common.Address, treeIDs []int, roots [][32]byte) error {
	chainID, err := sc.client.ChainID(sc.ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(sc.privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}

	nonce, err := sc.client.PendingNonceAt(sc.ctx, auth.From)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	gasPrice, err := sc.client.SuggestGasPrice(sc.ctx)
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %w", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)      // no ether being sent
	auth.GasLimit = uint64(3300000) // adjust as needed
	auth.GasPrice = gasPrice

	// Convert to big.Int
	treeIDsBig := make([]*big.Int, len(treeIDs))
	for i, id := range treeIDs {
		treeIDsBig[i] = big.NewInt(int64(id))
	}

	// Send the Merkle root to the smart contract
	tx, err := sc.contract.BatchUpdateTreeRoots(
		auth,
		issuers,
		treeIDsBig,
		roots,
	)
	if err != nil {
		return fmt.Errorf("failed to send root to contract: %w", err)
	}
	fmt.Printf("Transaction sent: %s\n", tx.Hash().Hex())

	// Wait for the transaction to be mined
	receipt, err := bind.WaitMined(sc.ctx, sc.client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for transaction to be mined: %w", err)
	}

	fmt.Printf("Transaction mined! Block number: %d\n", receipt.BlockNumber)

	return nil
}

func (sc *SmartContract) Verify(ctx context.Context, issuer common.Address, treeIndex int, leaf [32]byte, proof [][32]byte) (bool, error) {
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
