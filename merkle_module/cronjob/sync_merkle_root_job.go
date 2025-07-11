package cronjob

import (
	"context"
	"log"
	"merkle_module/domain/repo"
	"merkle_module/infra/model"
	credential "merkle_module/smartcontract"
	"merkle_module/utils"
	"sort"

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
	trees, err := j.repo.GetNodesToSync(j.ctx)
	if err != nil {
		return nil, err
	}

	var results []*RootResult
	for _, treeWithNodes := range trees {
		isValid, err := j.ValidateTree(treeWithNodes)
		if err != nil {
			log.Printf("Error validating tree ID %d: %v", treeWithNodes.Tree.ID, err)
			continue
		}
		if !isValid {
			log.Printf("Tree ID %d is not valid, skipping", treeWithNodes.Tree.ID)
			continue
		}
		var root []byte
		if len(treeWithNodes.Nodes) == 1 { // If there's only one node, use its data as the root
			root = treeWithNodes.Nodes[0].Data
		} else {
			// Build the merkle tree
			tree, err := mt.New(utils.GetTreeConfig(), utils.ToBlockDatas(treeWithNodes.Nodes))
			if err != nil {
				log.Printf("Error creating Merkle tree for Tree ID %d: %v", treeWithNodes.Tree.ID, err)
				continue
			}
			if tree == nil {
				log.Printf("Failed to create Merkle tree for Tree ID %d: tree is nil", treeWithNodes.Tree.ID)
				continue
			}
			root = tree.Root
		}

		// add the result
		results = append(results, &RootResult{
			Root:          root,
			TreeID:        treeWithNodes.Tree.ID,
			IssuerAddress: treeWithNodes.Tree.IssuerDID,
		})
	}

	return results, nil
}

func (j *SyncMerkleJob) ValidateTree(tree model.MerkleTreeWithNodes) (bool, error) {
	// Check if the tree has nodes
	if tree.Tree == nil || len(tree.Nodes) == 0 {
		log.Printf("Tree ID %d has no nodes, skipping validation", tree.Tree.ID)
		return false, nil
	}

	// check if nodes belong to the same tree
	sort.Slice(tree.Nodes, func(i, j int) bool {
		return tree.Nodes[i].NodeID < tree.Nodes[j].NodeID
	})
	for i := 1; i < len(tree.Nodes); i++ {
		if tree.Nodes[i].TreeID != tree.Nodes[i-1].TreeID {
			log.Printf("Nodes in tree ID %d do not belong to the same tree", tree.Tree.ID)
			return false, nil
		}
	}

	return true, nil
}
