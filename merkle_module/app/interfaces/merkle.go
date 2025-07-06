package interfaces

import "context"

type Merkle interface {
	AddLeaf(ctx context.Context, issuerDID string, data []byte) error
	GetProof(ctx context.Context, issuerDID string, data []byte) ([][]byte, error)
	VerifyProof(ctx context.Context, issuerDID string, data []byte, proof [][]byte) (bool, error)
}
