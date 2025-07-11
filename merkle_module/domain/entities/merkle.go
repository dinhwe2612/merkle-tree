package entities

type MerkleNode struct {
	ID     int    `json:"id"`
	TreeID int    `json:"tree_id"`
	NodeID int    `json:"node_id"`
	Data   []byte `json:"value"`
}

type MerkleTree struct {
	ID        int    `json:"id"`
	IssuerDID string `json:"issuer_did"`
	NeedSync  bool   `json:"need_sync"`
}

func (node *MerkleNode) Serialize() ([]byte, error) {
	// Serialize the node data to bytes
	return node.Data, nil
}
