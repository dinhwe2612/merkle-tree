package storage

import (
	"context"
	"fmt"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"
	"sync"
)

type MerklesInMemory struct {
	treeMap map[string]*merkletree.MerkleTree // issuerDID -> MerkleTree
	lock    sync.Mutex
}

func NewMerklesInMemory() repo.MerklesCache {
	return &MerklesInMemory{
		treeMap: make(map[string]*merkletree.MerkleTree),
		lock:    sync.Mutex{},
	}
}

func (m *MerklesInMemory) GetTree(ctx context.Context, issuerDID string) (*merkletree.MerkleTree, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	tree, exists := m.treeMap[issuerDID]
	if !exists {
		return nil, nil // Tree not found
	}

	// If the tree is full, remove it from the cache and return nil
	if tree.IsFull() {
		delete(m.treeMap, issuerDID)
		return nil, nil
	}

	return tree, nil
}

func (m *MerklesInMemory) BuildTree(ctx context.Context, issuerDID string, treeID int, nodes []string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Convert nodes to byte slices
	byteNodes := make([][]byte, len(nodes))
	for i, node := range nodes {
		byteNodes[i] = []byte(node)
	}

	// Create a new Merkle tree
	tree, err := merkletree.NewMerkleTree(byteNodes, treeID)
	if err != nil {
		return fmt.Errorf("failed to create new merkle tree: %w", err)
	}
	if tree == nil {
		return fmt.Errorf("failed to create new merkle tree: tree is nil")
	}

	m.treeMap[issuerDID] = tree

	return nil
}

func (m *MerklesInMemory) GetTreeByIssuerDIDAndData(ctx context.Context, issuerDID string, data []byte) *merkletree.MerkleTree {
	m.lock.Lock()
	defer m.lock.Unlock()

	// Get the tree for the given issuer DID
	tree, exists := m.treeMap[issuerDID]
	if !exists {
		return nil
	}

	// Check if the data exists in the tree
	if !tree.Contains(data) {
		return nil
	}

	return tree
}
