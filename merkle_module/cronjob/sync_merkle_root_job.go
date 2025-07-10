package cronjob

import (
	"context"
	"log"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"
	credential "merkle_module/smartcontract"
	"merkle_module/utils"

	"github.com/ethereum/go-ethereum/common"
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
	Root   []byte
	TreeID int
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
		log.Printf("Tree ID: %d, Merkle Root: %x", result.TreeID, result.Root)
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
		issuers = append(issuers, common.HexToAddress("0xYourIssuerAddress"))
		roots = append(roots, root)
		treeIDs = append(treeIDs, result.TreeID)
	}
	if err := j.contract.SendRoot(issuers, treeIDs, roots); err != nil {
		log.Printf("Error sending Merkle roots to smart contract: %v", err)
		return
	}

	log.Println("Merkle roots successfully sent to smart contract")
}

func (j *SyncMerkleJob) getRootResults() ([]RootResult, error) {
	results, err := j.repo.GetTreesWithNodesForSync(j.ctx)
	if err != nil {
		return nil, err
	}

	var rootResults []RootResult
	for _, tree := range results {
		// check node ids must fill the range from 1 to NodeCount
		if tree.Tree.NodeCount != len(tree.Nodes) {
			log.Printf("Node count mismatch for Tree ID %d: expected %d, got %d", tree.Tree.ID, tree.Tree.NodeCount, len(tree.Nodes))
			log.Printf("Skipping Tree ID %d due to invalid node IDs", tree.Tree.ID)
			continue
		}
		nodeMap := make(map[int]bool)
		flag := false
		for _, node := range tree.Nodes {
			if node.NodeID <= 0 || node.NodeID > tree.Tree.NodeCount {
				log.Printf("Invalid node ID %d for Tree ID %d: must be between 1 and %d", node.NodeID, tree.Tree.ID, tree.Tree.NodeCount)
				flag = true
			}
			if nodeMap[node.NodeID] {
				log.Printf("Duplicate node ID %d found for Tree ID %d", node.NodeID, tree.Tree.ID)
				flag = true
			}
			nodeMap[node.NodeID] = true
		}
		if flag {
			log.Printf("Skipping Tree ID %d due to invalid node IDs", tree.Tree.ID)
			continue
		}

		// build the Merkle tree
		tree, err := merkletree.NewMerkleTree(utils.NodesToBytes(tree.Nodes), tree.Tree.ID)
		if err != nil {
			log.Printf("Error creating Merkle tree for Tree ID %d: %v", tree.GetTreeID(), err)
			continue
		}

		// get the Merkle root
		root := tree.GetMerkleRoot()
		if root == nil {
			log.Printf("Error getting Merkle root for Tree ID %d: root is nil", tree.GetTreeID())
			continue
		}

		// append the result
		rootResults = append(rootResults, RootResult{
			Root:   root,
			TreeID: tree.GetTreeID(),
		})
	}

	return rootResults, nil
}
