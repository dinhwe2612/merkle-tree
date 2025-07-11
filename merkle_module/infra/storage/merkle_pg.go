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

func (m *MerklePostgres) AddNode(ctx context.Context, issuerDID string, data []byte) error {
	// Insert the new node
	_, err := m.db.ExecContext(ctx, `
	INSERT INTO merkle_nodes (issuer_did, data)
	VALUES ($1, $2)
	`, issuerDID, data)
	if err != nil {
		return fmt.Errorf("failed to insert new node: %w", err)
	}

	return nil
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

func (m *MerklePostgres) GetNodesToSync(ctx context.Context) ([][]*entities.MerkleNode, error) {
	// Get all nodes that need to be synced
	// Set 1 contains nodes from trees with less than MAX_LEAFS
	// Set 2 contains nodes from the tree with ID -1 (which is a special case for syncing)
	// Prevent case when there are no new nodes to sync by checking if set2 exists
	rows, err := m.db.QueryContext(ctx, `
	WITH set1 AS (
		SELECT mn.id, mn.tree_id, mn.node_id, mn.data, mn.issuer_did
		FROM merkle_nodes mn
		JOIN merkle_trees mt ON mn.tree_id = mt.tree_id
		WHERE mt.node_count < $1
	),
	set2 AS (
		SELECT id, tree_id, node_id, data, issuer_did
		FROM merkle_nodes
		WHERE tree_id = -1
	)
	SELECT * FROM (
		SELECT * FROM set1
		UNION
		SELECT * FROM set2
	) AS combined
	WHERE EXISTS (SELECT 1 FROM set2)
	ORDER BY issuer_did, tree_id DESC, node_id;
	`, utils.MAX_LEAFS)
	if err != nil {
		return nil, fmt.Errorf("failed to query merkle nodes: %w", err)
	}
	defer rows.Close()

	// group nodes by issuer DID
	var nodesOfIssuers [][]*entities.MerkleNode
	for rows.Next() {
		var node entities.MerkleNode
		if err := rows.Scan(&node.ID, &node.TreeID, &node.NodeID, &node.Data, &node.IssuerDID); err != nil {
			return nil, fmt.Errorf("failed to scan node: %w", err)
		}

		// If we are still in the same issuer, append to the current list
		if len(nodesOfIssuers) > 0 && nodesOfIssuers[len(nodesOfIssuers)-1][0].IssuerDID == node.IssuerDID {
			nodesOfIssuers[len(nodesOfIssuers)-1] = append(nodesOfIssuers[len(nodesOfIssuers)-1], &node)
		} else {
			// Otherwise, start a new list for this issuer
			nodesOfIssuers = append(nodesOfIssuers, []*entities.MerkleNode{&node})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return nodesOfIssuers, nil
}

func (m *MerklePostgres) GetTreeIDByData(ctx context.Context, hashValue []byte) (int, error) {
	var treeID int
	err := m.db.QueryRowContext(ctx, `
	SELECT tree_id
	FROM merkle_nodes
	WHERE data = $1
	LIMIT 1
	`, hashValue).Scan(&treeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no tree found for data %x", hashValue)
		}
		return 0, fmt.Errorf("failed to get tree ID by data: %w", err)
	}
	return treeID, nil
}

func (m *MerklePostgres) UpdateNodes(ctx context.Context, nodes []*entities.MerkleNode) error {
	if len(nodes) == 0 {
		return nil
	}

	stmt, err := m.db.PrepareContext(ctx, `
	UPDATE merkle_nodes
	SET tree_id = $1, node_id = $2
	WHERE id = $3
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	for _, node := range nodes {
		_, err := stmt.ExecContext(ctx, node.TreeID, node.NodeID, node.ID)
		if err != nil {
			return fmt.Errorf("failed to update node ID %d: %w", node.ID, err)
		}
	}

	return nil
}

func (m *MerklePostgres) UpdateTrees(ctx context.Context, trees []*entities.MerkleTree) error {
	if len(trees) == 0 {
		return nil
	}

	stmt, err := m.db.PrepareContext(ctx, `
	INSERT INTO merkle_trees (tree_id, issuer_did, node_count)
	VALUES ($1, $2, $3)
	ON CONFLICT (tree_id, issuer_did) DO UPDATE
	SET node_count = EXCLUDED.node_count
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement for trees: %w", err)
	}
	defer stmt.Close()

	for _, tree := range trees {
		_, err := stmt.ExecContext(ctx, tree.TreeID, tree.IssuerDID, tree.NodeCount)
		if err != nil {
			return fmt.Errorf("failed to upsert tree %d: %w", tree.TreeID, err)
		}
	}

	return nil
}
