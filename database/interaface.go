package database

type MerkleRepo interface {
	GetValues(nodeIDs []int, issuerDID string, treeID int) ([]string, error)
	UpdateValues(nodeIDs []int, values []string, issuerDID string, treeID int) error
	GetLeaf(issuerDID string, treeID int) ([]string, error)
	GetProof(issuerDID string, treeID int, data []byte) ([]string, error)
	GetNumberOfLeafs(issuerDID string, treeID int) (int, error)
	// AddLeaf(issuerDID string, treeID int, data []byte) (string, error)
}
