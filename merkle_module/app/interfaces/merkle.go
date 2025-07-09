package interfaces

import (
	"context"
	"merkle_module/domain/entities"
)

type Merkle interface {
	AddLeaf(ctx context.Context, issuerDID string, hashValue []byte) (*entities.MerkleNode, error)
	// This function is used to get proof for the tree in database
	GetProof(ctx context.Context, treeID, nodeID int) ([][]byte, error)
	GetRoot(ctx context.Context, treeID int) ([]byte, error)
	// This function is used to get the proof that has been synced
	GetSyncedProof(ctx context.Context, treeID, nodeID int) ([][]byte, error)
	// This function is used to get the root that has been synced
	GetSyncedRoot(ctx context.Context, treeID int) ([]byte, error)
}
