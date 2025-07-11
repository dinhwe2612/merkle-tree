package interfaces

import (
	"context"
	"merkle_module/domain/entities"
)

type Merkle interface {
	AddLeaf(ctx context.Context, issuerDID string, hashValue []byte) (*entities.MerkleNode, error)
	GetProof(ctx context.Context, treeID int, hashValue []byte) ([][]byte, error)
	GetRoot(ctx context.Context, treeID int) ([]byte, error)
}
