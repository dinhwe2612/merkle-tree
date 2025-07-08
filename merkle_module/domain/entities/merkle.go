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
	NodeCount int    `json:"node_count"`
	NeedSync  bool   `json:"need_sync"`
}
