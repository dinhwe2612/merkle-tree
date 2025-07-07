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
	repo               repo.Merkle
	cacheTrees         *lru.Cache[int, *merkletree.MerkleTree] // cache the Merkle trees
	cacheActiveTreeIDs *lru.Cache[string, int]                 // cache the active Merkle tree IDs
}

func NewMerkleService(repo repo.Merkle, cacheTrees *lru.Cache[int, *merkletree.MerkleTree], cacheActiveTreeIDs *lru.Cache[string, int]) interfaces.Merkle {
	return &MerkleService{repo: repo, cacheTrees: cacheTrees, cacheActiveTreeIDs: cacheActiveTreeIDs}
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
	tree, exists := s.cacheTrees.Get(treeID)

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
	_ = s.cacheTrees.Add(treeID, tree)

	return tree, nil
}

// helper function to get the active tree ID for an issuer DID
// , it also check if the active tree ID still correct, if not, return not found
func (s *MerkleService) getActiveTreeID(ctx context.Context, issuerDID string) (int, bool) {
	// Check if the active tree ID is cached
	activeTreeID, exists := s.cacheActiveTreeIDs.Get(issuerDID)
	if exists {
		// check if the tree is full or not
		tree, exists := s.cacheTrees.Get(activeTreeID)
		if exists && tree != nil {
			if !tree.IsFull() {
				return activeTreeID, true // return the cached active tree ID if the tree is not full
			}
		}
	}

	return 0, false // return 0 and false if the active tree ID is not cached or the tree is full
}

func (s *MerkleService) AddLeaf(ctx context.Context, issuerDID string, data []byte) (*entities.MerkleNode, error) {
	// Check if the active tree ID of the issuer DID is cached
	activeTreeID, exists := s.getActiveTreeID(ctx, issuerDID)
	if !exists {
		// If not cached, get the active tree for inserting
		nodes, err := s.repo.GetActiveTreeForInserting(ctx, issuerDID)
		if err != nil {
			return nil, fmt.Errorf("failed to get active tree for inserting: %w", err)
		}

		// Create a new Merkle tree
		tree, err := merkletree.NewMerkleTree(nodes.Nodes, nodes.TreeID)
		if err != nil {
			return nil, fmt.Errorf("failed to create new merkle tree: %w", err)
		}

		activeTreeID = nodes.TreeID
		s.cacheTrees.Add(activeTreeID, tree)
		s.cacheActiveTreeIDs.Add(issuerDID, activeTreeID)
	}

	// Get the tree by active tree ID
	tree, err := s.getTree(ctx, activeTreeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tree: %w", err)
	}

	// Add the leaf to the tree
	tree.AddLeaf(data)

	// Get the node ID of the added leaf
	nodeID, err := tree.GetLastNodeID()
	if err != nil {
		return nil, fmt.Errorf("failed to get last node ID: %w", err)
	}

	// Add the node to the database
	node, err := s.repo.AddNode(ctx, tree.GetTreeID(), nodeID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to add node: %w", err)
	}

	return node, nil
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
