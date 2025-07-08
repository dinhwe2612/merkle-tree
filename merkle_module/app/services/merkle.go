package services

import (
	"context"
	"fmt"
	"merkle_module/app/interfaces"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"
	"sync"

	"github.com/ethereum/go-ethereum/common/lru"
)

type MerkleService struct {
	repo               repo.Merkle
	cacheTrees         *lru.Cache[int, *merkletree.MerkleTree] // cache the Merkle trees
	cacheActiveTreeIDs *lru.Cache[string, int]                 // cache the active Merkle tree IDs
}

var muxtexes sync.Map // map to hold mutexes for each issuer DID

func getMutex(issuerDID string) *sync.Mutex {
	actual, _ := muxtexes.LoadOrStore(issuerDID, &sync.Mutex{})
	return actual.(*sync.Mutex)
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

// helper function to get the active tree
// for tree full, it remove the current active tree from the cache and return false
func (s *MerkleService) getActiveTree(ctx context.Context, issuerDID string) (*merkletree.MerkleTree, bool) {
	// Check if the active tree ID of the issuer DID is cached
	activeTreeID, exists := s.cacheActiveTreeIDs.Get(issuerDID)
	if !exists {
		return nil, false
	}

	// Get the tree from the cache
	tree, err := s.getTree(ctx, activeTreeID)
	if err != nil {
		return nil, false
	}

	// If the tree is full, remove it from the cache and return false
	if tree.IsFull() {
		s.cacheActiveTreeIDs.Remove(issuerDID)
		return nil, false
	}

	return tree, true
}

// helper function to get the active tree for inserting, return tree and boolean indicating if we need to load it from the database
func (s *MerkleService) getActiveTreeForInserting(ctx context.Context, issuerDID string) (*merkletree.MerkleTree, bool, error) {
	// Check if the active tree ID of the issuer DID is cached
	tree, exists := s.getActiveTree(ctx, issuerDID)

	if !exists {
		// If not cached, get the active tree for inserting
		activeTree, err := s.repo.GetActiveTreeForInserting(ctx, issuerDID)
		if err != nil {
			return nil, true, fmt.Errorf("failed to get active tree for inserting: %w", err)
		}

		// Create a new Merkle tree
		tree, err = merkletree.NewMerkleTree(activeTree.Nodes, activeTree.TreeID)
		if err != nil {
			return nil, true, fmt.Errorf("failed to create new merkle tree: %w", err)
		}

		s.cacheTrees.Add(tree.GetTreeID(), tree)
		s.cacheActiveTreeIDs.Add(issuerDID, tree.GetTreeID())

		return tree, true, nil
	}

	return tree, false, nil
}

func (s *MerkleService) AddLeaf(ctx context.Context, issuerDID string, data []byte) (*entities.MerkleNode, error) {
	mutex := getMutex(issuerDID)
	mutex.Lock()

	tree, needToLoad, err := s.getActiveTreeForInserting(ctx, issuerDID)
	if err != nil {
		mutex.Unlock()
		return nil, fmt.Errorf("failed to get active tree for inserting: %w", err)
	}

	// Add the leaf to the tree
	nodeID := tree.AddLeaf(data)
	if nodeID < 0 {
		mutex.Unlock()
		return nil, fmt.Errorf("failed to add leaf to tree: node ID is negative")
	}

	mutex.Unlock()

	// Add the node to the database
	var node *entities.MerkleNode
	if needToLoad {
		node, err = s.repo.AddNode(ctx, tree.GetTreeID(), nodeID, data)
	} else {
		node, err = s.repo.AddNodeAndIncrementNodeCount(ctx, tree.GetTreeID(), nodeID, data)
	}
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
