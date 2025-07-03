package database

type Merkle interface {
	GetValues(nodeIDs []int, issuerDID string, treeID string) ([]string, error)
	UpdateValues(nodeIDs []int, values []string, issuerDID string, treeID string) error
	GetLeaf(issuerDID string, treeID string) ([]string, error)
	GetProof(issuerDID string, treeID string, data string) (string, error)
	GetNumberOfLeafs(issuerDID string, treeID string) (int, error)
}
