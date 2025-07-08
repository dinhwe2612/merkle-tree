package repo

import (
	"context"
	"merkle_module/domain/entities"
	"merkle_module/infra/model"
)

type Merkle interface {
	// Get all nodes belonging to a specific tree ID
	GetNodesByTreeID(ctx context.Context, treeID int) ([][]byte, error)
	AddNode(ctx context.Context, treeID int, nodeID int, data []byte) (*entities.MerkleNode, error)
	// Retrieve the active tree and reserve an empty node for inserting a new one
	GetActiveTreeForInserting(ctx context.Context, issuerDID string) (*model.ActiveTree, error)
	// Add a new node to the tree and increment the node count
	AddNodeAndIncrementNodeCount(ctx context.Context, treeID int, nodeID int, data []byte) (*entities.MerkleNode, error)
	// Get the nodes of trees that need to be synced
	GetTreesWithNodesForSync(ctx context.Context) ([]*model.MerkleTreeWithNodes, error)
}
