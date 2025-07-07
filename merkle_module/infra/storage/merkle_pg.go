package storage

import (
	"context"
	"database/sql"
	"fmt"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	"merkle_module/utils"
)

type MerklePostgres struct {
	db *sql.DB
}

func NewMerklePostgres(db *sql.DB) repo.Merkle {
	return &MerklePostgres{db: db}
}

func (m *MerklePostgres) GetNodesByTreeID(ctx context.Context, treeID int) ([]string, error) {
	rows, err := m.db.QueryContext(ctx, `
	SELECT node_id, value
	FROM merkle_nodes
	WHERE tree_id = $1
	ORDER BY node_id
	`, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle nodes: %w", err)
	}
	defer rows.Close()

	var nodes []string
	for rows.Next() {
		var nodeID int
		var nodeData []byte
		if err := rows.Scan(&nodeID, &nodeData); err != nil {
			return nil, fmt.Errorf("failed to scan merkle node: %w", err)
		}
		nodes = append(nodes, fmt.Sprintf("%d:%s", nodeID, nodeData))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return nodes, nil
}

func (m *MerklePostgres) GetActiveTreeID(ctx context.Context, issuerDID string) (int, error) {
	var treeID int
	err := m.db.QueryRowContext(ctx, `
	SELECT id
	FROM merkle_trees
	WHERE issuer_did = $1 AND node_count < $2
	`, issuerDID, utils.MAX_LEAFS).Scan(&treeID)

	// If no tree is found, create a new one
	if err == sql.ErrNoRows {
		err = m.db.QueryRowContext(ctx, `
		INSERT INTO merkle_trees (issuer_did, node_count)
		VALUES ($1, 0)
		RETURNING id
		`, issuerDID).Scan(&treeID)
		if err != nil {
			return 0, fmt.Errorf("failed to create new merkle tree: %w", err)
		}

		return treeID, nil
	}

	if err != nil {
		return 0, fmt.Errorf("failed to query active tree ID: %w", err)
	}

	return treeID, nil
}

func (m *MerklePostgres) GetNodesByIssuerDIDAndData(ctx context.Context, issuerDID string, data string) ([]string, error) {
	// Get tree ID of the data belong to the issuer DID
	var treeID int
	err := m.db.QueryRowContext(ctx, `
	SELECT tree_id
	FROM merkle_trees
	WHERE issuer_did = $1 AND value = $2
	`, issuerDID, data).Scan(&treeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no tree found for issuer DID: %s with data: %s", issuerDID, data)
		}
		return nil, fmt.Errorf("failed to query tree ID: %w", err)
	}

	rows, err := m.db.QueryContext(ctx, `
	SELECT node_id, value
	FROM merkle_nodes
	WHERE tree_id = $1 AND value = $2
	ORDER BY node_id
	`, treeID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle nodes: %w", err)
	}
	defer rows.Close()

	var nodes []string
	for rows.Next() {
		var nodeID int
		var nodeData []byte
		if err := rows.Scan(&nodeID, &nodeData); err != nil {
			return nil, fmt.Errorf("failed to scan merkle node: %w", err)
		}
		nodes = append(nodes, fmt.Sprintf("%d:%s", nodeID, nodeData))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for issuer DID: %s with data: %s", issuerDID, data)
	}

	return nodes, nil
}

func (m *MerklePostgres) AddNode(ctx context.Context, issuerDID string, treeID int, nodeID int, data string) (*entities.MerkleNode, error) {
	// Insert the new node into the database
	_, err := m.db.ExecContext(ctx, `
	INSERT INTO merkle_nodes (tree_id, node_id, value)
	VALUES ($1, $2, $3)
	`, treeID, nodeID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to insert merkle node: %w", err)
	}

	// Return the new node
	return &entities.MerkleNode{
		TreeID: treeID,
		NodeID: nodeID,
	}, nil
}

func (m *MerklePostgres) GetTreeIDByIssuerDID(ctx context.Context, issuerDID string) (int, error) {
	var treeID int
	err := m.db.QueryRowContext(ctx, `
	SELECT id
	FROM merkle_trees
	WHERE issuer_did = $1
	`, issuerDID).Scan(&treeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no tree found for issuer DID: %s", issuerDID)
		}
		return 0, fmt.Errorf("failed to query tree ID: %w", err)
	}
	return treeID, nil
}
