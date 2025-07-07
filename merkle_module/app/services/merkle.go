package services

import (
	"context"
	"fmt"
	"merkle_module/app/interfaces"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"
)

type MerkleService struct {
	repo  repo.Merkle
	cache repo.MerklesCache // store active trees
}

func NewMerkleService(repo repo.Merkle, cache repo.MerklesCache) interfaces.Merkle {
	return &MerkleService{repo: repo, cache: cache}
}

// helper function to get active tree from cache or build it from database
func (s *MerkleService) getActiveTree(ctx context.Context, issuerDID string) (*merkletree.MerkleTree, error) {
	// Get the tree from the cache
	tree, err := s.cache.GetTree(ctx, issuerDID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree from cache: %w", err)
	}

	// If the tree is not found in the cache, create a new one
	if tree == nil {
		// Get active tree ID from database
		treeID, err := s.repo.GetActiveTreeID(ctx, issuerDID)
		if err != nil {
			return nil, fmt.Errorf("failed to get active tree ID: %w", err)
		}

		// Get nodes of tree id from database
		nodes, err := s.repo.GetNodesByTreeID(ctx, treeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get nodes by tree ID: %w", err)
		}

		// Build the tree from cache
		err = s.cache.BuildTree(ctx, issuerDID, treeID, nodes)
		if err != nil {
			return nil, fmt.Errorf("failed to create new merkle tree: %w", err)
		}

		// Get tree again
		tree, err = s.cache.GetTree(ctx, issuerDID)
		if err != nil {
			return nil, fmt.Errorf("failed to get tree from cache after building: %w", err)
		}
	}

	return tree, nil
}

// get tree of data from cache or build it from database
func (s *MerkleService) getTreeByIssuerDIDAndData(ctx context.Context, issuerDID string, data []byte) (*merkletree.MerkleTree, error) {
	// Get the tree from the cache
	tree := s.cache.GetTreeByIssuerDIDAndData(ctx, issuerDID, data)

	if tree != nil {
		return tree, nil
	}

	// If not in cache, build the tree from the database
	nodes, err := s.repo.GetNodesByIssuerDIDAndData(ctx, issuerDID, string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes by issuer DID and data: %w", err)
	}

	// Convert nodes to byte slices
	var byteNodes [][]byte
	for _, node := range nodes {
		byteNodes = append(byteNodes, []byte(node))
	}

	// Create a new Merkle tree
	tree, err = merkletree.NewMerkleTree(byteNodes, 0) // 0 is a placeholder for treeID, as we don't need it here
	if err != nil {
		return nil, fmt.Errorf("failed to create new merkle tree: %w", err)
	}

	if tree == nil {
		return nil, fmt.Errorf("failed to create new merkle tree: tree is nil")
	}

	return tree, nil
}

func (s *MerkleService) AddLeaf(ctx context.Context, issuerDID string, data []byte) (*entities.MerkleNode, error) {
	tree, err := s.getActiveTree(ctx, issuerDID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active tree: %w", err)
	}

	// Save to tree
	tree.AddLeaf(data)

	// Get the node ID from the cache
	nodeID, err := tree.GetLastNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to get last node ID from tree: %w", err)
	}

	// Save to database
	merkleNode, err := s.repo.AddNode(ctx, issuerDID, tree.GetTreeID(), nodeID, string(data))
	if err != nil {
		return nil, fmt.Errorf("failed to add node to database: %w", err)
	}

	return merkleNode, nil
}

func (s *MerkleService) GetProof(ctx context.Context, issuerDID string, data []byte) ([][]byte, error) {
	// Get the tree by issuer DID and data
	tree, err := s.getTreeByIssuerDIDAndData(ctx, issuerDID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree by issuer DID and data: %w", err)
	}

	// Get the proof for the data
	proof, err := tree.GetProof(data)

	return proof, nil
}

func (s *MerkleService) GetRoot(ctx context.Context, issuerDID string, data []byte) ([]byte, error) {
	// Get the tree by issuer DID and data
	tree, err := s.getTreeByIssuerDIDAndData(ctx, issuerDID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree by issuer DID and data: %w", err)
	}

	return tree.GetMerkleRoot(), nil
}
