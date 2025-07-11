package cronjob

import (
	"context"
	"fmt"
	"log"
	"merkle_module/domain/entities"
	"merkle_module/domain/repo"
	credential "merkle_module/smartcontract"
	"merkle_module/utils"

	"github.com/ethereum/go-ethereum/common"
	mt "github.com/txaty/go-merkletree"
)

type SyncMerkleJob struct {
	ctx      context.Context
	repo     repo.Merkle
	contract *credential.SmartContract
}

func NewSyncMerkleJob(ctx context.Context, repo repo.Merkle, contract *credential.SmartContract) *SyncMerkleJob {
	return &SyncMerkleJob{
		ctx:      ctx,
		repo:     repo,
		contract: contract,
	}
}

type RootResult struct {
	Root          []byte
	TreeID        int
	IssuerAddress string
}

// Run executes the job to sync the Merkle root, implementing the Job interface.
func (j *SyncMerkleJob) Run() {
	// Get the results of the Merkle root sync
	results, err := j.getRootResults()
	if err != nil {
		log.Printf("Error getting root results: %v", err)
		return
	}

	// If no results, exit early
	if results == nil || len(results) == 0 {
		log.Println("No results to sync")
		return
	}

	// Print the results
	for _, result := range results {
		log.Printf("Tree ID: %d, Issuer Address: %s, Merkle Root: %x", result.TreeID, result.IssuerAddress, result.Root)
	}

	// Send to smart contract
	if j.contract == nil {
		log.Println("Smart contract is not initialized, skipping sending Merkle roots")
		return
	}
	var issuers []common.Address
	var roots [][32]byte
	var treeIDs []int
	for _, result := range results {
		// Convert the Merkle root to a 32-byte array
		if len(result.Root) != 32 {
			log.Printf("Invalid Merkle root length for Tree ID %d: expected 32 bytes, got %d bytes", result.TreeID, len(result.Root))
			continue
		}
		var root [32]byte
		copy(root[:], result.Root)

		// Append the issuer address and root
		issuers = append(issuers, common.HexToAddress(result.IssuerAddress))
		roots = append(roots, root)
		treeIDs = append(treeIDs, result.TreeID)
	}
	if err := j.contract.SendRoot(issuers, treeIDs, roots); err != nil {
		log.Printf("Error sending Merkle roots to smart contract: %v", err)
		return
	}

	log.Println("Merkle roots successfully sent to smart contract")
}

func (j *SyncMerkleJob) getRootResults() ([]*RootResult, error) {
	nodesOfIssuers, err := j.repo.GetNodesToSync(j.ctx)
	if err != nil {
		return nil, err
	}

	// group nodes of issuerDID into trees following by utils.MAX_LEAFS
	var updatedNodes []*entities.MerkleNode
	var results []*RootResult
	var updatedTree []*entities.MerkleTree
	for _, nodes := range nodesOfIssuers {
		if len(nodes) == 0 {
			continue
		}
		fmt.Printf("Processing %d nodes for issuer %s\n", len(nodes), nodes[0].IssuerDID)
		// Index the nodeID and treeID again
		currentTreeID := nodes[0].TreeID
		if currentTreeID <= 0 {
			currentTreeID = 1 // Start from tree ID 1 if not set
		}
		var datas [][]byte
		for _, node := range nodes {
			if len(datas) >= utils.MAX_LEAFS {
				// create a merkle tree from the current data
				updatedTree = append(updatedTree, &entities.MerkleTree{
					IssuerDID: node.IssuerDID,
					NodeCount: len(datas),
					TreeID:    currentTreeID,
				})
				root, err := j.getRoot(datas)
				if err != nil {
					return nil, err
				}
				results = append(results, &RootResult{
					Root:          root,
					TreeID:        currentTreeID,
					IssuerAddress: node.IssuerDID,
				})
				// go to the next tree
				currentTreeID++
				datas = nil
			}
			// Update the node with the new tree ID and node ID
			node.TreeID = currentTreeID
			node.NodeID = len(datas) + 1
			updatedNodes = append(updatedNodes, node)
			datas = append(datas, node.Data)
		}
		if len(datas) > 0 {
			// create a merkle tree from the remaining data
			updatedTree = append(updatedTree, &entities.MerkleTree{
				IssuerDID: nodes[0].IssuerDID,
				NodeCount: len(datas),
				TreeID:    currentTreeID,
			})

			// Get the Merkle root for the current tree
			root, err := j.getRoot(datas)
			if err != nil {
				return nil, err
			}

			results = append(results, &RootResult{
				Root:          root,
				TreeID:        currentTreeID,
				IssuerAddress: nodes[0].IssuerDID,
			})
		}
	}

	// Update the nodes in the database
	if err := j.repo.UpdateNodes(j.ctx, updatedNodes); err != nil {
		return nil, err
	}

	// Update the trees
	if err := j.repo.UpdateTrees(j.ctx, updatedTree); err != nil {
		return nil, err
	}

	return results, nil
}

func (j *SyncMerkleJob) getRoot(datas [][]byte) ([]byte, error) {
	if len(datas) == 0 {
		return nil, fmt.Errorf("no data provided to create Merkle root")
	}

	if len(datas) == 1 {
		return datas[0], nil
	}

	// Create a Merkle tree from the data
	tree, err := mt.New(utils.GetTreeConfig(), utils.ToBlockDataFromByteArray(datas))
	if err != nil {
		return nil, err
	}

	// Return the Merkle root
	return tree.Root, nil
}
