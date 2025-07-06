package storage

import (
	"context"
	"database/sql"
	"fmt"
	"merkle_module/domain/repo"
	"merkle_module/utils"
)

type MerklePostgres struct {
	db *sql.DB
}

func NewMerklePostgres(db *sql.DB) repo.Merkle {
	return &MerklePostgres{db: db}
}

func (m *MerklePostgres) GetTreeIDByIssuerDIDAndValue(ctx context.Context, issuerDID string, value string) (int, error) {
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

func (m *MerklePostgres) GetNodesByTreeID(ctx context.Context, issuerDID string, treeID int) ([]string, error) {
	rows, err := m.db.QueryContext(ctx, `
	SELECT nodes.node_id, nodes.value 
	FROM merkle_nodes nodes
	JOIN merkle_trees trees ON nodes.tree_id = trees.id
	WHERE trees.issuer_did = $1 AND nodes.tree_id = $2
	ORDER BY nodes.node_id
	`, issuerDID, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query nodes by tree ID: %w", err)
	}
	defer rows.Close()

	var nodes []string
	lastID := 1
	for rows.Next() {
		var nodeID int
		var value string
		if err := rows.Scan(&nodeID, &value); err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}
		// append the empty value between the gap
		for i := lastID; i < nodeID; i++ {
			nodes = append(nodes, "")
		}
		lastID = nodeID + 1
		nodes = append(nodes, value)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return nodes, nil
}

func (m *MerklePostgres) AddNode(ctx context.Context, issuerDID string, treeID int, value string, nodeID int) error {
	// insert the new node
	_, err := m.db.ExecContext(ctx, `
	INSERT INTO merkle_nodes (tree_id, node_id, value)
	VALUES ($1, $2, $3)
	`, treeID, nodeID, value)
	if err != nil {
		return fmt.Errorf("failed to insert node: %w", err)
	}

	return nil
}

// gets or creates a tree, returns the new node ID of the tree, and increases the node count
// if the current tree is full (node_count >= MAX_LEAFS), it creates a new tree
func (m *MerklePostgres) GetNewNodeIDAndIncreaseCount(ctx context.Context, issuerDID string) (int, int, error) {
	// Start transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to begin transaction: %w", err)
	}

	// find an existing tree not full and increase node_count in one query
	var treeID int
	var newNodeID int
	err = tx.QueryRowContext(ctx, `
	UPDATE merkle_trees
	SET node_count = node_count + 1
	WHERE issuer_did = $1 AND node_count < $2
	RETURNING id, node_count
	`, issuerDID, utils.MAX_LEAFS).Scan(&treeID, &newNodeID)

	if err == nil {
		if err = tx.Commit(); err != nil {
			tx.Rollback()
			return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
		}

		return treeID, newNodeID, nil
	}

	if err != sql.ErrNoRows {
		return 0, 0, fmt.Errorf("failed to query existing tree: %w", err)
	}

	// create a new tree
	err = tx.QueryRowContext(ctx, `
	INSERT INTO merkle_trees (issuer_did, node_count)
	VALUES ($1, 1)
	RETURNING id
	`, issuerDID).Scan(&treeID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create new tree: %w", err)
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return 0, 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return treeID, 1, nil
}
