package entities

type MerkleNode struct {
	ID        int    `json:"id"`
	TreeID    int    `json:"tree_id"`
	NodeID    int    `json:"node_id"`
	Data      []byte `json:"value"`
	IssuerDID string `json:"issuer_did"`
}

type MerkleTree struct {
	ID        int    `json:"id"`
	IssuerDID string `json:"issuer_did"`
	NodeCount int    `json:"node_count"`
	TreeID    int    `json:"tree_id"`
}

func (node *MerkleNode) Serialize() ([]byte, error) {
	// Serialize the node data to bytes
	return node.Data, nil
}
