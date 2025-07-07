package storage

import (
	"context"
	"database/sql"
	"fmt"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	"merkle_module/infra/model"
	"merkle_module/utils"
)

type MerklePostgres struct {
	db *sql.DB
}

func NewMerklePostgres(db *sql.DB) repo.Merkle {
	return &MerklePostgres{db: db}
}

func (m *MerklePostgres) GetNodesByTreeID(ctx context.Context, treeID int) ([][]byte, error) {
	rows, err := m.db.QueryContext(ctx, `
	SELECT node_id, data
	FROM merkle_nodes
	WHERE tree_id = $1
	ORDER BY node_id
	`, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle nodes: %w", err)
	}
	defer rows.Close()

	var datas [][]byte
	for rows.Next() {
		var nodeID int
		var nodeData []byte
		if err := rows.Scan(&nodeID, &nodeData); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		datas = append(datas, nodeData)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return datas, nil
}

func (m *MerklePostgres) GetActiveTreeID(ctx context.Context, issuerDID string) (int, error) {
	// Begin a transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get the active tree ID and put
	var treeID int
	err = tx.QueryRowContext(ctx, `
	SELECT id
	FROM merkle_trees
	WHERE issuer_did = $1 AND node_count < $2
	FOR UPDATE
	`, issuerDID, utils.MAX_LEAFS).Scan(&treeID)

	// If not found, create a new one
	if err == sql.ErrNoRows {
		err = tx.QueryRowContext(ctx, `
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

	// commit the transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return treeID, nil
}

func (m *MerklePostgres) AddNode(ctx context.Context, treeID int, nodeID int, data []byte) (*entities.MerkleNode, error) {
	// Insert the new node into the database
	_, err := m.db.ExecContext(ctx, `
	INSERT INTO merkle_nodes (tree_id, node_id, data)
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

func (m *MerklePostgres) GetActiveTreeForInserting(ctx context.Context, issuerDID string) (*model.ActiveTree, error) {
	// Begin a transaction
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Get the active tree ID for the issuer DID
	var treeID int
	err = tx.QueryRowContext(ctx, `
	SELECT id 
	FROM merkle_trees 
	WHERE issuer_did = $1 AND node_count < $2 
	FOR UPDATE
	`, issuerDID, utils.MAX_LEAFS).Scan(&treeID)
	if err == sql.ErrNoRows {
		// If not found, create a new one
		err = tx.QueryRowContext(ctx, `
		INSERT INTO merkle_trees (issuer_did, node_count) 
		VALUES ($1, 1) 
		RETURNING id
		`, issuerDID).Scan(&treeID)
		if err != nil {
			return nil, fmt.Errorf("failed to create new merkle tree: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get active tree id: %w", err)
	} else {
		// If found, increase node_count by 1 to reserve a slot
		_, err = tx.ExecContext(ctx, `
		UPDATE merkle_trees 
		SET node_count = node_count + 1 
		WHERE id = $1
		`, treeID)
		if err != nil {
			return nil, fmt.Errorf("failed to update node_count: %w", err)
		}
	}

	// Get the nodes for the active tree
	rows, err := tx.QueryContext(ctx, `
	SELECT data 
	FROM merkle_nodes 
	WHERE tree_id = $1 
	ORDER BY node_id
	`, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes by tree ID: %w", err)
	}
	defer rows.Close()
	var nodes [][]byte
	for rows.Next() {
		var nodeData []byte
		if err := rows.Scan(&nodeData); err != nil {
			return nil, fmt.Errorf("failed to scan node data: %w", err)
		}
		nodes = append(nodes, nodeData)
	}

	return &model.ActiveTree{
		TreeID:    treeID,
		IssuerDID: issuerDID,
		NodeCount: len(nodes) + 1, // Already reserved a slot for the next node
		Nodes:     nodes,
	}, nil
}
