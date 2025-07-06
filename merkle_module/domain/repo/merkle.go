package repo

import (
	"context"
	"merkle_module/domain/entities"
)

type Merkle interface {
	GetTreeIDByIssuerDIDAndData(ctx context.Context, issuerDID string, data string) (int, error)
	GetNodesByTreeID(ctx context.Context, treeID int) ([]string, error)
	AddNode(ctx context.Context, issuerDID string, data string) (*entities.MerkleNode, error)
}

type MerklesCache interface {
	HasData(ctx context.Context, issuerDID string, data []byte) (bool, error)
	HasTree(ctx context.Context, treeID int) (bool, error)
	AddNode(ctx context.Context, issuerDID string, treeID int, data []byte) error
	LoadTree(ctx context.Context, issuerDID string, treeID int, datas [][]byte) error
	GetProof(ctx context.Context, issuerDID string, data []byte) ([][]byte, error)
	GetRoot(ctx context.Context, issuerDID string, data []byte) ([]byte, error)
}
