package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	// make a copy of the data
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	// Insert the new node into the database
	_, err := m.db.ExecContext(ctx, `
	INSERT INTO merkle_nodes (tree_id, node_id, data)
	VALUES ($1, $2, $3)
	`, treeID, nodeID, dataCopy)
	if err != nil {
		return nil, fmt.Errorf("failed to insert merkle node: %w", err)
	}

	// Return the new node
	return &entities.MerkleNode{
		TreeID: treeID,
		NodeID: nodeID,
		Data:   dataCopy,
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
	SET node_count = node_count + 1,
		need_sync = TRUE
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
	// Update and get the tree IDs that need to be synced
	var treeIDs []int64
	rows, err := m.db.QueryContext(ctx, `
	UPDATE merkle_trees
	SET need_sync = FALSE
	WHERE need_sync = TRUE
	RETURNING id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to update and get tree IDs for sync: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan tree ID: %w", err)
		}
		treeIDs = append(treeIDs, id)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error after scanning: %w", err)
	}

	if len(treeIDs) == 0 {
		return nil, nil // No trees to sync
	}

	// Get the nodes for the trees that need to be synced
	nodeRows, err := m.db.QueryContext(ctx, `
	SELECT mt.id AS tree_id, mn.node_id, mn.data
	FROM merkle_nodes mn
	JOIN merkle_trees mt ON mn.tree_id = mt.id
	WHERE mt.id = ANY($1) AND mn.node_id <= mt.node_count
	ORDER BY mt.id, mn.node_id
	`, pq.Array(treeIDs))
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle nodes for trees: %w", err)
	}
	defer nodeRows.Close()

	var currentTreeID int
	var currentTree *model.MerkleTreeWithNodes
	var result []*model.MerkleTreeWithNodes
	for nodeRows.Next() {
		var treeID, nodeID int
		var nodeData []byte
		if err := nodeRows.Scan(&treeID, &nodeID, &nodeData); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		log.Printf("Processing node ID %d for tree ID %d", nodeID, treeID)

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

	if err := nodeRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over nodeRows: %w", err)
	}

	// Update the node_count_sync for each tree using
	stmt, err := m.db.PrepareContext(ctx, `
	UPDATE merkle_trees
	SET node_count_sync = $1
	WHERE id = $2
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	for _, tree := range result {
		_, err := stmt.ExecContext(ctx, len(tree.Nodes), tree.Tree.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update node_count_sync for tree ID %d: %w", tree.Tree.ID, err)
		}
	}

	return result, nil
}

func (m *MerklePostgres) GetNodesSyncedByTreeID(ctx context.Context, treeID int) ([]*entities.MerkleNode, error) {
	// Get the nodes synced by tree ID
	rows, err := m.db.QueryContext(ctx, `
	SELECT node_id, data
	FROM merkle_nodes mn
	JOIN merkle_trees mt ON mn.tree_id = mt.id
	WHERE mn.tree_id = $1 and mn.node_id <= mt.node_count_sync
	ORDER BY mn.node_id
	`, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle nodes by tree ID: %w", err)
	}
	defer rows.Close()

	var nodes []*entities.MerkleNode
	for rows.Next() {
		var nodeID int
		var nodeData []byte
		if err := rows.Scan(&nodeID, &nodeData); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		nodes = append(nodes, &entities.MerkleNode{
			TreeID: treeID,
			NodeID: nodeID,
			Data:   nodeData,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return nodes, nil
}
