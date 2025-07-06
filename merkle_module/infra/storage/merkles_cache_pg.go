package storage

import (
	"context"
	"fmt"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"
	"merkle_module/utils"
	"sync"
)

type MerklesInMemory struct {
	treeMap map[int]*merkletree.MerkleTree
	leafMap map[string]int // issuer_did | data -> treeID
	lock    sync.Mutex
}

func NewMerklesInMemory() repo.MerklesCache {
	return &MerklesInMemory{
		treeMap: make(map[int]*merkletree.MerkleTree),
		leafMap: make(map[string]int),
	}
}

func (m *MerklesInMemory) HasData(ctx context.Context, issuerDID string, data []byte) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := fmt.Sprintf("%s|%s", issuerDID, utils.Hash(data))
	_, exists := m.leafMap[key]
	return exists, nil
}

func (m *MerklesInMemory) HasTree(ctx context.Context, treeID int) (bool, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	_, exists := m.treeMap[treeID]
	return exists, nil
}

func (m *MerklesInMemory) AddNode(ctx context.Context, issuerDID string, treeID int, data []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := fmt.Sprintf("%s|%s", issuerDID, utils.Hash(data))
	if _, exists := m.leafMap[key]; exists {
		return fmt.Errorf("data already exists in cache")
	}

	tree, exists := m.treeMap[treeID]
	if !exists {
		return fmt.Errorf("tree with ID %d does not exist in cache", treeID)
	}

	tree.AddLeaf(data)
	m.leafMap[key] = treeID

	return nil
}

func (m *MerklesInMemory) LoadTree(ctx context.Context, issuerDID string, treeID int, datas [][]byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	tree, err := merkletree.NewMerkleTree(datas)
	if err != nil {
		return fmt.Errorf("failed to create new Merkle tree: %w", err)
	}
	m.treeMap[treeID] = tree

	for _, data := range datas {
		key := fmt.Sprintf("%s|%s", issuerDID, utils.Hash(data))
		m.leafMap[key] = treeID
	}

	return nil
}

func (m *MerklesInMemory) GetProof(ctx context.Context, issuerDID string, data []byte) ([][]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := fmt.Sprintf("%s|%s", issuerDID, utils.Hash(data))
	treeID, exists := m.leafMap[key]
	if !exists {
		return nil, fmt.Errorf("data not found in cache")
	}

	tree, exists := m.treeMap[treeID]
	if !exists {
		return nil, fmt.Errorf("tree with ID %d does not exist in cache", treeID)
	}

	proof, err := tree.GetProof(data)
	if err != nil {
		return nil, fmt.Errorf("failed to get proof: %w", err)
	}

	return proof, nil
}

func (m *MerklesInMemory) GetRoot(ctx context.Context, issuerDID string, data []byte) ([]byte, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	key := fmt.Sprintf("%s|%s", issuerDID, utils.Hash(data))
	treeID, exists := m.leafMap[key]
	if !exists {
		return nil, fmt.Errorf("data not found in cache")
	}

	tree, exists := m.treeMap[treeID]
	if !exists {
		return nil, fmt.Errorf("tree with ID %d does not exist in cache", treeID)
	}

	return tree.GetMerkleRoot(), nil
}
