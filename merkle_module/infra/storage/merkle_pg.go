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

func (m *MerklePostgres) GetTreeIDByIssuerDIDAndData(ctx context.Context, issuerDID string, value string) (int, error) {
	var treeID int
	err := m.db.QueryRowContext(ctx, `
	SELECT nodes.tree_id
	FROM merkle_nodes nodes
	JOIN merkle_trees trees ON nodes.tree_id = trees.id
	WHERE trees.issuer_did = $1 AND nodes.value = $2
	`, issuerDID, value).Scan(&treeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("tree not found for value: %s", value)
		}
		return 0, fmt.Errorf("failed to get tree ID: %w", err)
	}

	return treeID, nil
}

func (m *MerklePostgres) GetNodesByTreeID(ctx context.Context, treeID int) ([]string, error) {
	rows, err := m.db.QueryContext(ctx, `
	SELECT nodes.value 
	FROM merkle_nodes nodes
	WHERE nodes.tree_id = $1
	ORDER BY nodes.node_id
	`, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query nodes by tree ID: %w", err)
	}
	defer rows.Close()

	var nodes []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		nodes = append(nodes, value)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return nodes, nil
}

func (m *MerklePostgres) AddNode(ctx context.Context, issuerDID string, value string) (*entities.MerkleNode, error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if the data already exists for this issuer
	if exists, err := m.checkDataExists(ctx, tx, issuerDID, value); err != nil {
		return nil, err
	} else if exists {
		return nil, fmt.Errorf("data already exists for issuer %s: %s", issuerDID, value)
	}

	// Find an existing tree not full and increase node_count
	var treeID int
	var newNodeID int
	err = tx.QueryRowContext(ctx, `
	UPDATE merkle_trees 
	SET node_count = node_count + 1
	WHERE issuer_did = $1 AND node_count < $2
	RETURNING id, node_count
	`, issuerDID, utils.MAX_LEAFS).Scan(&treeID, &newNodeID)

	if err == sql.ErrNoRows {
		// Create a new tree
		err = tx.QueryRowContext(ctx, `
		INSERT INTO merkle_trees (issuer_did, node_count)
		VALUES ($1, 1)
		RETURNING id
		`, issuerDID).Scan(&treeID)
		if err != nil {
			return nil, fmt.Errorf("failed to create new tree: %w", err)
		}
		newNodeID = 1
	} else if err != nil {
		return nil, fmt.Errorf("failed to update existing tree: %w", err)
	}

	// Insert the new node
	var nodeID int
	err = tx.QueryRowContext(ctx, `
	INSERT INTO merkle_nodes (tree_id, node_id, value)
	VALUES ($1, $2, $3)
	RETURNING id
	`, treeID, newNodeID, value).Scan(&nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to insert node: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &entities.MerkleNode{
		ID:     nodeID,
		TreeID: treeID,
		NodeID: newNodeID,
		Value:  value,
	}, nil
}

func (m *MerklePostgres) checkDataExists(ctx context.Context, tx *sql.Tx, issuerDID string, value string) (bool, error) {
	var exists bool
	err := tx.QueryRowContext(ctx, `
	SELECT EXISTS (
		SELECT 1 
		FROM merkle_nodes nodes
		JOIN merkle_trees trees ON nodes.tree_id = trees.id
		WHERE trees.issuer_did = $1 AND nodes.value = $2
	)
	`, issuerDID, value).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if data exists: %w", err)
	}
	return exists, nil
}
