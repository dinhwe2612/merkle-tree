package entities

import "time"

// MerkleNode represents a node in the Merkle tree stored in the database
type MerkleNode struct {
	ID     int    `json:"id"`
	TreeID int    `json:"tree_id"`
	NodeID int    `json:"node_id"`
	Value  string `json:"value"`
	// Optional fields that might be useful
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

// MerkleTree represents a Merkle tree stored in the database
type MerkleTree struct {
	ID        int    `json:"id"`
	IssuerDID string `json:"issuer_did"`
	NodeCount int    `json:"node_count"`
	// Optional fields that might be useful
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
