package cronjob

import (
	"context"
	"log"
	"merkle_module/domain/repo"
	"merkle_module/merkletree"
	"merkle_module/utils"
)

type SyncMerkleJob struct {
	ctx  context.Context
	repo repo.Merkle
}

func NewSyncMerkleJob(ctx context.Context, repo repo.Merkle) *SyncMerkleJob {
	return &SyncMerkleJob{
		ctx:  ctx,
		repo: repo,
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
	// for _, result := range results {
	// 	// Convert root to byte32
	// 	root32 := utils.ToByte32(result.Root)
	// 	if err != nil {
	// 		log.Printf("Error converting root to byte32 for Tree ID %d: %v", result.TreeID, err)
	// 		continue
	// 	}

	// 	// Send the Merkle root to the smart contract
	// 	sendsmartcontract.SendRootToContract(root32)
	// }
}

func (j *SyncMerkleJob) getRootResults() ([]RootResult, error) {
	results, err := j.repo.GetTreesWithNodesForSync(j.ctx)
	if err != nil {
		return nil, err
	}

	var rootResults []RootResult
	for _, tree := range results {
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
		rootResults = append(rootResults, RootResult{
			Root:   root,
			TreeID: tree.GetTreeID(),
		})
	}
	return rootResults, nil
}
