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

	datas := make([][]byte, 0)
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
	// fmt.Printf("Retrieved %d nodes for tree ID %d\n", len(datas), treeID)
	return datas, nil
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

	var treeID int
	var nodeID int
	err = tx.QueryRowContext(ctx, `
	SELECT id, node_count 
	FROM merkle_trees 
	WHERE issuer_did = $1 AND node_count < $2 
	FOR UPDATE
	`, issuerDID, utils.MAX_LEAFS).Scan(&treeID, &nodeID)

	if err == sql.ErrNoRows {
		// If not found, create a new one with node_count = 1
		err = tx.QueryRowContext(ctx, `
		INSERT INTO merkle_trees (issuer_did, node_count) 
		VALUES ($1, 1) 
		RETURNING id, node_count
		`, issuerDID).Scan(&treeID, &nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to create new merkle tree: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to get active tree id: %w", err)
	} else {
		// If found, increase node_count by 1 and get the new value
		err = tx.QueryRowContext(ctx, `
		UPDATE merkle_trees 
		SET node_count = node_count + 1,
			need_sync = TRUE
		WHERE id = $1 
		RETURNING node_count
		`, treeID).Scan(&nodeID)
		if err != nil {
			return nil, fmt.Errorf("failed to update and return node_count: %w", err)
		}
	}

	fmt.Printf("Retrieved nodes for tree ID %d, node ID %d\n", treeID, nodeID)

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
		NodeCount: nodeID, // Already reserved a slot for the next node
		Nodes:     nodes,
	}, nil
}

func (m *MerklePostgres) AddNodeAndIncrementNodeCount(ctx context.Context, treeID int, nodeID int, data []byte) (*entities.MerkleNode, error) {
	// Insert the new node into the database
	_, err := m.db.ExecContext(ctx, `
	INSERT INTO merkle_nodes (tree_id, node_id, data)
	VALUES ($1, $2, $3)
	`, treeID, nodeID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to insert merkle node: %w", err)
	}

	// Increment the node count for the tree with lock
	_, err = m.db.ExecContext(ctx, `
	UPDATE merkle_trees
	SET node_count = node_count + 1
	WHERE id = $1
	`, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to increment node count: %w", err)
	}

	// Return the new node
	return &entities.MerkleNode{
		TreeID: treeID,
		NodeID: nodeID,
	}, nil
}

func (m *MerklePostgres) GetTreesWithNodesForSync(ctx context.Context) ([]*model.MerkleTreeWithNodes, error) {
	// Query the trees and their nodes that need sync
	rows, err := m.db.QueryContext(ctx, `
	WITH updated AS (
		UPDATE merkle_trees
		SET need_sync = FALSE
		WHERE need_sync = TRUE
		RETURNING id
	)
	SELECT t.id, n.node_id, n.data
	FROM updated t
	JOIN merkle_nodes n ON t.id = n.tree_id
	ORDER BY t.id, n.node_id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle trees and nodes: %w", err)
	}
	defer rows.Close()

	var currentTreeID int
	var currentTree *model.MerkleTreeWithNodes
	var result []*model.MerkleTreeWithNodes
	for rows.Next() {
		var treeID, nodeID int
		var nodeData []byte
		if err := rows.Scan(&treeID, &nodeID, &nodeData); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if currentTreeID != treeID {
			if currentTree != nil {
				result = append(result, currentTree)
			}
			currentTreeID = treeID
			currentTree = &model.MerkleTreeWithNodes{
				Tree: &entities.MerkleTree{
					ID:        treeID,
					NodeCount: 0,
				},
				Nodes: []*entities.MerkleNode{},
			}
		}
		currentTree.Nodes = append(currentTree.Nodes, &entities.MerkleNode{
			TreeID: treeID,
			NodeID: nodeID,
			Data:   nodeData,
		})
		currentTree.Tree.NodeCount++
	}

	if currentTree != nil {
		result = append(result, currentTree)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return result, nil
}
