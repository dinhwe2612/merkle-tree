package storage

import (
	"context"
	"database/sql"
	"fmt"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	"merkle_module/infra/model"
	"merkle_module/utils"

	"github.com/lib/pq"
)

type MerklePostgres struct {
	db *sql.DB
}

func NewMerklePostgres(db *sql.DB) repo.Merkle {
	return &MerklePostgres{db: db}
}

func (m *MerklePostgres) AddNode(ctx context.Context, issuerDID string, data []byte) (node *entities.MerkleNode, err error) {
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				err = fmt.Errorf("failed to rollback transaction: %w, original error: %v", rbErr, err)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				err = fmt.Errorf("failed to commit transaction: %w", commitErr)
			}
		}
	}()

	// Acquire advisory lock for issuerDID
	_, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtext($1))`, issuerDID)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire advisory lock: %w", err)
	}

	// Check if there's an existing data for the issuerDID
	var existingCount int
	err = tx.QueryRowContext(ctx, `
	SELECT COUNT(*)
	FROM merkle_nodes
	WHERE tree_id IN (
		SELECT id
		FROM merkle_trees
		WHERE issuer_did = $1
	) AND data = $2
	`, issuerDID, data).Scan(&existingCount)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing data: %w", err)
	}
	if existingCount > 0 {
		return nil, fmt.Errorf("data already exists for issuerDID %s", issuerDID)
	}

	var treeID, nodeCount int
	err = tx.QueryRowContext(ctx, `
	UPDATE merkle_trees
	SET node_count = node_count + 1, need_sync = true
	WHERE issuer_did = $1 AND node_count < $2
	RETURNING id, node_count
	`, issuerDID, utils.MAX_LEAFS).Scan(&treeID, &nodeCount)

	if err == sql.ErrNoRows {
		// No available tree — create new one
		err = tx.QueryRowContext(ctx, `
		INSERT INTO merkle_trees (issuer_did, node_count, need_sync)
		VALUES ($1, 1, true)
		RETURNING id, node_count
		`, issuerDID).Scan(&treeID, &nodeCount)
		if err != nil {
			return nil, fmt.Errorf("failed to create new tree: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to update existing tree: %w", err)
	}

	// Insert node
	_, err = tx.ExecContext(ctx, `
	INSERT INTO merkle_nodes (tree_id, data)
	VALUES ($1, $2)
	`, treeID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to insert node: %w", err)
	}

	return &entities.MerkleNode{
		TreeID: treeID,
		Data:   data,
	}, nil
}

func (m *MerklePostgres) GetNodesToBuildTree(ctx context.Context, treeID int) ([][]byte, error) {
	rows, err := m.db.QueryContext(ctx, `
		SELECT data
		FROM merkle_nodes
		WHERE tree_id = $1 AND node_id <> 0
		ORDER BY node_id
	`, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes by tree ID %d: %w", treeID, err)
	}
	defer rows.Close()

	var nodes [][]byte
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed to scan node data: %w", err)
		}
		nodes = append(nodes, data)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for tree ID %d", treeID)
	}

	return nodes, nil
}

func (m *MerklePostgres) GetNodesToSync(ctx context.Context) ([]model.MerkleTreeWithNodes, error) {
	// Step 1: Fetch trees that need syncing
	rows, err := m.db.QueryContext(ctx, `
		UPDATE merkle_trees
		SET need_sync = false
		WHERE need_sync = true
		RETURNING id, issuer_did
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get trees to sync: %w", err)
	}
	defer rows.Close()

	var treeIDs []int
	var results []model.MerkleTreeWithNodes
	treeIndex := make(map[int]int) // tree_id → index in results

	for rows.Next() {
		tree := &entities.MerkleTree{}
		if err := rows.Scan(&tree.ID, &tree.IssuerDID); err != nil {
			return nil, fmt.Errorf("failed to scan tree: %w", err)
		}
		treeIndex[tree.ID] = len(results)
		results = append(results, model.MerkleTreeWithNodes{Tree: tree})
		treeIDs = append(treeIDs, tree.ID)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating result rows: %w", err)
	}
	if len(results) == 0 {
		return nil, nil // No trees to sync
	}

	// Step 2: Fetch nodes for the selected trees
	rows, err = m.db.QueryContext(ctx, `
		SELECT n.tree_id, n.id, n.data
		FROM merkle_nodes n
		WHERE n.tree_id = ANY($1)
		ORDER BY n.tree_id, n.id
	`, pq.Array(treeIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes for trees: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var treeID, id int
		var data []byte
		if err := rows.Scan(&treeID, &id, &data); err != nil {
			return nil, fmt.Errorf("failed to scan node data: %w", err)
		}
		idx, ok := treeIndex[treeID]
		if !ok {
			return nil, fmt.Errorf("unexpected tree_id %d in nodes", treeID)
		}
		results[idx].Nodes = append(results[idx].Nodes, &entities.MerkleNode{
			ID:     id,
			TreeID: treeID,
			Data:   data,
			NodeID: len(results[idx].Nodes) + 1,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating node rows: %w", err)
	}

	// Step 3: Assign and update node_id based on order within each tree
	stmt, err := m.db.PrepareContext(ctx, `
		UPDATE merkle_nodes
		SET node_id = $1
		WHERE id = $2
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	for _, treeWithNodes := range results {
		for i, node := range treeWithNodes.Nodes {
			node.NodeID = i // node_id is set to index (0-based)
			_, err := stmt.ExecContext(ctx, i, node.NodeID)
			if err != nil {
				return nil, fmt.Errorf("failed to update node_id for node %d: %w", node.ID, err)
			}
		}
	}

	return results, nil
}
