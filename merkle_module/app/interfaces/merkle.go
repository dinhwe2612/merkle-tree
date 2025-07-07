package interfaces

import (
	"context"
	"merkle_module/domain/entities"
)

type Merkle interface {
	AddLeaf(ctx context.Context, issuerDID string, data []byte) (*entities.MerkleNode, error)
	GetProof(ctx context.Context, issuerDID string, data []byte) ([][]byte, error)
	GetRoot(ctx context.Context, issuerDID string, data []byte) ([]byte, error)
}
