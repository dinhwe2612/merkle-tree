package repo

import "context"

type Merkle interface {
	GetTreeIDByIssuerDIDAndValue(ctx context.Context, issuerDID string, value string) (int, error)
	GetNodesByTreeID(ctx context.Context, issuerDID string, treeID int) ([]string, error)
	GetNewNodeIDAndIncreaseCount(ctx context.Context, issuerDID string) (int, int, error)
	AddNode(ctx context.Context, issuerDID string, treeID int, value string, nodeID int) error
}
