package database

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type MerklePostgres struct {
	db *sql.DB
}

func NewMerklePostgres(db *sql.DB) MerkleRepo {
	return &MerklePostgres{db: db}
}

func (r *MerklePostgres) GetValues(nodeIDs []int, issuerDID string, treeID int) ([]string, error) {
	if len(nodeIDs) == 0 {
		return []string{}, nil
	}

	rows, err := r.db.Query(`
	SELECT node_id, value 
	FROM merkle_node 
	WHERE node_id = ANY($1) AND tree_id = $2
	`, pq.Array(nodeIDs), treeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	valueMap := make(map[int]string)
	for rows.Next() {
		var nodeID int
		var value string
		if err := rows.Scan(&nodeID, &value); err != nil {
			return nil, err
		}
		valueMap[nodeID] = value
	}

	values := make([]string, len(nodeIDs))
	for i, id := range nodeIDs {
		values[i] = valueMap[id] // If not found, will be empty string
	}

	return values, nil
}

func (r *MerklePostgres) UpdateValues(nodeIDs []int, values []string, issuerDID string, treeID int) error {
	if len(nodeIDs) != len(values) {
		return fmt.Errorf("nodeIDs and values length mismatch: %d vs %d", len(nodeIDs), len(values))
	}
	if len(nodeIDs) == 0 {
		return nil
	}

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
	UPDATE merkle_node 
	SET value = $1 
	WHERE node_id = $2 AND tree_id = $3`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for i := range nodeIDs {
		_, err = stmt.Exec(values[i], nodeIDs[i], treeID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (r *MerklePostgres) GetLeaf(issuerDID string, treeID int) ([]string, error) {
	rows, err := r.db.Query(`
	SELECT value 
	FROM merkle_node 
	WHERE tree_id = $1 
	ORDER BY node_id
	`, treeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leafs []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		leafs = append(leafs, value)
	}

	return leafs, nil
}

func (r *MerklePostgres) GetProof(issuerDID string, treeID int, data []byte) ([]string, error) {
	// hash the data
	dataHash := Hash(data)

	// get the node ID of the data
	var nodeID int
	err := r.db.QueryRow(`
	SELECT node_id 
	FROM merkle_node
	WHERE value = $1 AND tree_id = $2
	`, dataHash, treeID).Scan(&nodeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("data not found in tree")
		}
		return nil, err
	}

	// calculate the list of sibling ids
	var siblingIDs []int
	for nodeID > 1 {
		siblingID := nodeID ^ 1 // Get the sibling ID
		siblingIDs = append(siblingIDs, siblingID)
		nodeID >>= 1 // Move up to the parent node
	}

	// get the values of the sibling nodes, order by node_id to maintain the proof order from leaf to root
	rows, err := r.db.Query(`
	SELECT value 
	FROM merkle_node 
	WHERE node_id = ANY($1) 
		and tree_id = $2
	ORDER BY node_id
	`, pq.Array(siblingIDs), treeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var proof []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		proof = append(proof, value)
	}

	return proof, nil
}

func (r *MerklePostgres) GetNumberOfLeafs(issuerDID string, treeID int) (int, error) {
	var count int
	err := r.db.QueryRow(`
	SELECT tree_count 
	FROM merkle_tree 
	WHERE tree_id = $1
	`, treeID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No tree found, return 0
		}
		return 0, err
	}

	return count, nil
}

// func (r *MerklePostgres) AddLeaf(issuerDID string, treeID int, data []byte) (string, error) {
// 	// hash the data
// 	dataHash := Hash(data)

// 	// get the left most nodeID that is not used
// 	var numLeafs int
// 	err := r.db.QueryRow(`
// 	SELECT tree_count
// 	FROM merkle_tree
// 	WHERE tree_id = $1
// 	`, treeID).Scan(&numLeafs)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return "", fmt.Errorf("tree not found")
// 		}
// 		return "", err
// 	}
// 	nodeID := MAX_LEAFS + numLeafs

// 	// calculate the list of sibling and parent ids
// 	siblingIDs := []int{}
// 	parentIDs := []int{}
// 	for nodeID > 1 {
// 		siblingID := nodeID ^ 1 // Get the sibling ID
// 		siblingIDs = append(siblingIDs, siblingID)
// 		parentID := nodeID >> 1 // Get the parent ID
// 		parentIDs = append(parentIDs, parentID)
// 		nodeID = parentID // Move up to the parent node
// 	}

// 	// get the values of the sibling and parent nodes
// 	rows, err := r.db.Query(`
// 	SELECT value
// 	FROM merkle_node
// 	WHERE node_id = ANY($1) AND tree_id = $2
// 	`, pq.Array(siblingIDs), treeID)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer rows.Close()
// 	var siblingValues []string
// 	for rows.Next() {
// 		var value string
// 		if err := rows.Scan(&value); err != nil {
// 			return "", err
// 		}
// 		siblingValues = append(siblingValues, value)
// 	}

// 	// update the sibling nodes with the new data
// }
