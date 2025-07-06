package services

import (
	"context"
	"fmt"
	"merkle_module/app/interfaces"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"
)

type MerkleService struct {
	repo repo.Merkle
}

func NewMerkleService(repo repo.Merkle) interfaces.Merkle {
	return &MerkleService{repo: repo}
}

func (s *MerkleService) AddLeaf(ctx context.Context, issuerDID string, data []byte) error {
	treeID, nodeID, err := s.repo.GetNextNodeIDAndIncreaseCount(ctx, issuerDID)
	if err != nil {
		return fmt.Errorf("failed to get next node ID: %w", err)
	}

	if err := s.repo.AddNode(ctx, issuerDID, treeID, string(data), nodeID); err != nil {
		return fmt.Errorf("failed to add node: %w", err)
	}

	return nil
}

func (s *MerkleService) GetProof(ctx context.Context, issuerDID string, data []byte) ([][]byte, error) {
	treeID, err := s.repo.GetTreeIDByValue(ctx, issuerDID, string(data))
	if err != nil {
		return nil, err
	}

	// get the leaf nodes to build the Merkle tree
	nodes, err := s.repo.GetNodesByTreeID(ctx, issuerDID, treeID)
	if err != nil {
		return nil, err
	}

	// convert []string to [][]byte
	nodeBytes := make([][]byte, len(nodes))
	for i, node := range nodes {
		nodeBytes[i] = []byte(node)
	}

	// build the tree
	tree, err := merkletree.NewMerkleTree(nodeBytes)
	if err != nil {
		return nil, err
	}

	// get the proof for the original data
	proof, err := tree.GetProof(data)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

func (s *MerkleService) VerifyProof(ctx context.Context, issuerDID string, data []byte, proof [][]byte) (bool, error) {
	// Find tree ID of the data
	treeID, err := s.repo.GetTreeIDByValue(ctx, issuerDID, string(data))
	if err != nil {
		return false, fmt.Errorf("failed to get tree ID: %w", err)
	}

	nodes, err := s.repo.GetNodesByTreeID(ctx, issuerDID, treeID)
	if err != nil {
		return false, fmt.Errorf("failed to get nodes: %w", err)
	}

	// convert []string to [][]byte
	nodeBytes := make([][]byte, len(nodes))
	for i, node := range nodes {
		nodeBytes[i] = []byte(node)
	}

	// build the tree
	tree, err := merkletree.NewMerkleTree(nodeBytes)
	if err != nil {
		return false, fmt.Errorf("failed to create merkle tree: %w", err)
	}

	// verify the proof
	isValid := merkletree.Verify(proof, tree.GetMerkleRoot(), data)
	return isValid, nil
}
