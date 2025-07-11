package services

import (
	"context"
	"fmt"
	"merkle_module/app/interfaces"
	"merkle_module/domain/repo"
	"merkle_module/utils"

	mt "github.com/txaty/go-merkletree"
)

type MerkleService struct {
	repo repo.Merkle
}

func NewMerkleService(repo repo.Merkle) interfaces.Merkle {
	return &MerkleService{repo: repo}
}

func (s *MerkleService) AddLeaf(ctx context.Context, issuerDID string, data []byte) error {
	err := s.repo.AddNode(ctx, issuerDID, data)
	if err != nil {
		return fmt.Errorf("failed to add node: %w", err)
	}
	return nil
}

func (s *MerkleService) GetProof(ctx context.Context, data []byte) ([][]byte, error) {
	// Get tree ID of the data
	treeID, err := s.repo.GetTreeIDByData(ctx, data)

	// Load node from database
	nodes, err := s.repo.GetNodesToBuildTree(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree by ID %d: %w", treeID, err)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for tree ID %d", treeID)
	}

	if len(nodes) == 1 {
		// If there's only one node, return an empty proof
		return [][]byte{}, nil
	}

	// Build the Merkle tree
	tree, err := mt.New(utils.GetTreeConfig(), utils.ToBlockDataFromByteArray(nodes))
	if err != nil {
		return nil, fmt.Errorf("failed to create Merkle tree: %w", err)
	}

	if tree == nil {
		return nil, fmt.Errorf("failed to create Merkle tree: tree is nil")
	}

	proof, err := tree.Proof(utils.ToBlockData(data))
	if err != nil {
		return nil, fmt.Errorf("failed to get proof: %w", err)
	}

	return proof.Siblings, nil
}
