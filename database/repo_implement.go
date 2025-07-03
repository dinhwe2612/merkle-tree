package database

import (
	"database/sql"
	"fmt"
)

type RepoImpelement struct {
	db *sql.DB
}

func NewRepoImplement(db *sql.DB) *RepoImpelement {
	return &RepoImpelement{db: db}
}

func (r *RepoImpelement) GetValues(nodeIDs []int, issuerDID string, treeID string) ([]string, error) {
	query := "SELECT value FROM merkle_tree WHERE node_id IN (?) AND issuer_did = ? AND tree_id = ?"
	args := make([]interface{}, len(nodeIDs)+2)
	for i, id := range nodeIDs {
		args[i] = id
	}
	args[len(nodeIDs)] = issuerDID
	args[len(nodeIDs)+1] = treeID

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

func (r *RepoImpelement) UpdateValues(nodeIDs []int, values []string, issuerDID string, treeID string) error {
	if len(nodeIDs) != len(values) {
		return fmt.Errorf("nodeIDs and values length mismatch: %d vs %d", len(nodeIDs), len(values))
	}
	// update query (prevent for-loop)
	query := "UPDATE merkle_tree SET value = CASE node_id "
	args := make([]interface{}, 0, len(nodeIDs)*2+2)
	for i, id := range nodeIDs {
		query += "WHEN ? THEN ? "
		args = append(args, id, values[i])
	}
	query += "END WHERE node_id IN (?) AND issuer_did = ? AND tree_id = ?"
	args = append(args, nodeIDs, issuerDID, treeID)
	_, err := r.db.Exec(query, args...)
	return err
}
