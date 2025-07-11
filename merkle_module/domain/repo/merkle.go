package repo

import (
	"context"
	"merkle_module/domain/entities"
	"merkle_module/infra/model"
)

type Merkle interface {
	AddNode(ctx context.Context, issuerDID string, data []byte) (*entities.MerkleNode, error)
	GetNodesToBuildTree(ctx context.Context, treeID int) ([][]byte, error)
	// Get nodes of trees that need to be synced and index them by tree ID
	GetNodesToSync(ctx context.Context) ([]model.MerkleTreeWithNodes, error)
}
