package services

import (
	"context"
	"fmt"
	"merkle_module/app/interfaces"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"

	"github.com/ethereum/go-ethereum/common/lru"
)

type MerkleService struct {
	repo  repo.Merkle
	cache *lru.Cache[int, *merkletree.MerkleTree] // cache the Merkle trees
}

func NewMerkleService(repo repo.Merkle, cache *lru.Cache[int, *merkletree.MerkleTree]) interfaces.Merkle {
	return &MerkleService{repo: repo, cache: cache}
}

// helper function to build a new Merkle tree from the database
func (s *MerkleService) buildTree(ctx context.Context, treeID int) (*merkletree.MerkleTree, error) {
	// Get nodes of tree id from database
	nodes, err := s.repo.GetNodesByTreeID(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes by tree ID: %w", err)
	}

	// Create a new Merkle tree
	tree, err := merkletree.NewMerkleTree(nodes, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to create new merkle tree: %w", err)
	}

	// After creating the tree, it should be initialized
	if tree == nil {
		return nil, fmt.Errorf("failed to create new merkle tree: tree is nil")
	}

	return tree, nil
}

// helper function to get tree from cache or build it from database
func (s *MerkleService) getTree(ctx context.Context, treeID int) (*merkletree.MerkleTree, error) {
	// Get the tree from the cache
	tree, exists := s.cache.Get(treeID)

	// If the tree exists in the cache, return it
	if exists && tree != nil {
		return tree, nil
	}

	// If the tree is not found in the cache, create a new one
	// Build the tree from the database
	tree, err := s.buildTree(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to build tree: %w", err)
	}

	if tree == nil {
		return nil, fmt.Errorf("failed to build tree: tree is nil")
	}

	// set the tree in the cache
	_ = s.cache.Add(treeID, tree)

	return tree, nil
}

func (s *MerkleService) AddLeaf(ctx context.Context, issuerDID string, data []byte) (*entities.MerkleNode, error) {
	// Get the active tree ID for the issuer DID
	treeID, err := s.repo.GetActiveTreeID(ctx, issuerDID)
	if err != nil {
		return nil, fmt.Errorf("failed to get active tree ID: %w", err)
	}

	// Get the tree by tree ID
	tree, err := s.getTree(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	// Save to tree
	tree.AddLeaf(data)

	// Get the node ID from the cache
	nodeID, err := tree.GetLastNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to get last node ID from tree: %w", err)
	}

	// Save to database
	merkleNode, err := s.repo.AddNode(ctx, tree.GetTreeID(), nodeID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to add node to database: %w", err)
	}

	return merkleNode, nil
}

func (s *MerkleService) GetProof(ctx context.Context, treeID, nodeID int) ([][]byte, error) {
	// Get the tree by tree ID
	tree, err := s.getTree(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	proof, err := tree.GetProof(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get proof: %w", err)
	}

	return proof, nil
}

func (s *MerkleService) GetRoot(ctx context.Context, treeID int) ([]byte, error) {
	// Get the tree by tree ID
	tree, err := s.getTree(ctx, treeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	return tree.GetMerkleRoot(), nil
}
