package repo

import (
	"context"
	"merkle_module/domain/entities"
)

type Merkle interface {
	// Get all nodes belonging to a specific tree ID
	GetNodesByTreeID(ctx context.Context, treeID int) ([][]byte, error)
	// get the active tree ID for the given issuer DID
	// if no active tree is found, a new one is created and returned
	GetActiveTreeID(ctx context.Context, issuerDID string) (int, error)
	AddNode(ctx context.Context, treeID int, nodeID int, data []byte) (*entities.MerkleNode, error)
}
