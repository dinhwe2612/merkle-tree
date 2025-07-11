package model

import "merkle_module/domain/entities"

type MerkleTreeWithNodes struct {
	Tree  *entities.MerkleTree
	Nodes []*entities.MerkleNode
}
