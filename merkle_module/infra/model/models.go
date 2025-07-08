package model

import "merkle_module/domain/entities"

type ActiveTree struct {
	TreeID    int      `json:"tree_id"`
	IssuerDID string   `json:"issuer_did"`
	NodeCount int      `json:"node_count"`
	Nodes     [][]byte `json:"nodes"` // List of nodes in the tree
}

type MerkleTreeWithNodes struct {
	Tree  *entities.MerkleTree
	Nodes []*entities.MerkleNode
}
