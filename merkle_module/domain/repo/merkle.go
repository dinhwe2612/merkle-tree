package repo

import (
	"context"
	"merkle_module/domain/entities"
)

type Merkle interface {
	AddNode(ctx context.Context, issuerDID string, data []byte) error
	GetNodesToBuildTree(ctx context.Context, treeID int) ([][]byte, error)
	// Get nodes of trees that need to be synced and index them by tree ID
	GetNodesToSync(context.Context) ([][]*entities.MerkleNode, error)
	UpdateNodes(ctx context.Context, nodes []*entities.MerkleNode) error
	GetTreeIDByData(ctx context.Context, hashValue []byte) (int, error)
	UpdateTrees(ctx context.Context, trees []*entities.MerkleTree) error
}
