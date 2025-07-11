package services

import (
	"context"
	"fmt"
	"merkle_module/app/interfaces"
	"merkle_module/domain/entities"
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

func (s *MerkleService) AddLeaf(ctx context.Context, issuerDID string, data []byte) (*entities.MerkleNode, error) {
	node, err := s.repo.AddNode(ctx, issuerDID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to add node: %w", err)
	}
	return node, nil
}

func (s *MerkleService) GetProof(ctx context.Context, treeID int, data []byte) ([][]byte, error) {
	// Get the tree by tree ID
	tree, err := s.getTree(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	proof, err := tree.Proof(utils.ToBlockData(data))
	if err != nil {
		return nil, fmt.Errorf("failed to get proof: %w", err)
	}

	return proof.Siblings, nil
}

func (s *MerkleService) GetRoot(ctx context.Context, treeID int) ([]byte, error) {
	// Get the tree by tree ID
	tree, err := s.getTree(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	return tree.Root, nil
}

func (s *MerkleService) getTree(ctx context.Context, treeID int) (*mt.MerkleTree, error) {
	// Load node from database
	nodes, err := s.repo.GetNodesToBuildTree(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree by ID %d: %w", treeID, err)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("no nodes found for tree ID %d", treeID)
	}

	// Build the Merkle tree
	tree, err := mt.New(utils.GetTreeConfig(), utils.ToBlockDataFromByteArray(nodes))
	if err != nil {
		return nil, fmt.Errorf("failed to create Merkle tree: %w", err)
	}

	if tree == nil {
		return nil, fmt.Errorf("failed to create Merkle tree: tree is nil")
	}

	return tree, nil
}
