package repo

import (
	"context"
	"merkle_module/domain/entities"
	"merkle_module/merkletree"
)

type Merkle interface {
	// Get all nodes belonging to a specific tree ID
	GetNodesByTreeID(ctx context.Context, treeID int) ([]string, error)
	// get the active tree ID for the given issuer DID
	// if no active tree is found, a new one is created and returned
	GetActiveTreeID(ctx context.Context, issuerDID string) (int, error)
	GetNodesByIssuerDIDAndData(ctx context.Context, issuerDID, data string) ([]string, error)
	AddNode(ctx context.Context, issuerDID string, treeID int, nodeID int, data string) (*entities.MerkleNode, error)
}

type MerklesCache interface {
	GetTree(ctx context.Context, issuerDID string) (*merkletree.MerkleTree, error)
	BuildTree(ctx context.Context, issuerDID string, treeID int, nodes []string) error
	GetTreeByIssuerDIDAndData(ctx context.Context, issuerDID string, data []byte) *merkletree.MerkleTree
}
