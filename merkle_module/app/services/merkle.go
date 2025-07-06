package services

import (
	"context"
	"encoding/hex"
	"fmt"
	"merkle_module/app/interfaces"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"
)

type MerkleService struct {
	repo  repo.Merkle
	cache repo.MerklesCache
}

func NewMerkleService(repo repo.Merkle, cache repo.MerklesCache) interfaces.Merkle {
	return &MerkleService{repo: repo, cache: cache}
}

func (s *MerkleService) AddLeaf(ctx context.Context, issuerDID string, data []byte) (*entities.MerkleNode, error) {
	// Add leaf into database
	merkleNode, err := s.repo.AddNode(ctx, issuerDID, hex.EncodeToString(data))
	if err != nil {
		return nil, fmt.Errorf("failed to add node: %w", err)
	}

	// Check if the tree is loaded in cache
	treeExists, err := s.cache.HasTree(ctx, merkleNode.TreeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if tree exists in cache: %w", err)
	}

	if !treeExists {
		// Tree is not loaded in cache, no need to sync
		return merkleNode, nil
	}

	// Update cache with the new node
	if err := s.cache.AddNode(ctx, issuerDID, merkleNode.TreeID, data); err != nil {
		return nil, fmt.Errorf("failed to add node to cache: %w", err)
	}

	return merkleNode, nil
}

func (s *MerkleService) GetProof(ctx context.Context, issuerDID string, data []byte) ([][]byte, error) {
	// Sync the leaf and tree in cache
	if err := s.syncLeaf(ctx, issuerDID, data); err != nil {
		return nil, fmt.Errorf("failed to sync leaf: %w", err)
	}

	// Get proof from cache
	proof, err := s.cache.GetProof(ctx, issuerDID, data)
	if err != nil {
		return nil, fmt.Errorf("failed to get proof from cache: %w", err)
	}

	return proof, nil
}

func (s *MerkleService) VerifyProof(ctx context.Context, issuerDID string, data []byte, proof [][]byte) (bool, error) {
	rootHash, err := s.cache.GetRoot(ctx, issuerDID, data)
	if err != nil {
		return false, fmt.Errorf("failed to get root hash from cache: %w", err)
	}

	isValid := merkletree.Verify(proof, rootHash, data)

	return isValid, nil
}

func (s *MerkleService) syncTree(ctx context.Context, issuerDID string, treeID int) error {
	nodes, err := s.repo.GetNodesByTreeID(ctx, treeID)
	if err != nil {
		return fmt.Errorf("failed to get nodes by tree ID: %w", err)
	}

	// Convert to bytes
	var byteNodes [][]byte
	for _, node := range nodes {
		data, err := hex.DecodeString(node)
		if err != nil {
			return fmt.Errorf("failed to decode node value: %w", err)
		}
		byteNodes = append(byteNodes, data)
	}

	if err := s.cache.LoadTree(ctx, issuerDID, treeID, byteNodes); err != nil {
		return fmt.Errorf("failed to load tree into cache: %w", err)
	}

	return nil
}

func (s *MerkleService) syncLeaf(ctx context.Context, issuerDID string, data []byte) error {
	// Check if the data exists in cache
	exists, err := s.cache.HasData(ctx, issuerDID, data)
	if err != nil {
		return fmt.Errorf("failed to check if data exists in cache: %w", err)
	}

	if exists {
		return nil
	}

	// Get the tree ID from the database
	treeID, err := s.repo.GetTreeIDByIssuerDIDAndData(ctx, issuerDID, hex.EncodeToString(data))
	if err != nil {
		return fmt.Errorf("failed to get tree ID by issuer DID and data: %w", err)
	}

	// Sync the tree
	if err := s.syncTree(ctx, issuerDID, treeID); err != nil {
		return fmt.Errorf("failed to sync tree: %w", err)
	}

	return nil
}
