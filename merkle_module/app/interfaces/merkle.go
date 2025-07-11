package interfaces

import (
	"context"
)

type Merkle interface {
	AddLeaf(ctx context.Context, issuerDID string, hashValue []byte) error
	GetProof(ctx context.Context, hashValue []byte) ([][]byte, error)
}
