package cronjob

import (
	"context"
	"log"
	"merkle_module/domain/repo"
	credential "merkle_module/smartcontract"
)

type AsyncJob struct {
	ctx           context.Context
	jobManager    *JobManager
	repo          repo.Merkle
	smartContract *credential.SmartContract
}

func NewAsyncJob(ctx context.Context, repo repo.Merkle, smartContract *credential.SmartContract) *AsyncJob {
	jobManager := NewJobManager()
	return &AsyncJob{
		ctx:           ctx,
		jobManager:    jobManager,
		repo:          repo,
		smartContract: smartContract,
	}
}

func (aj *AsyncJob) Start() {
	// Add a job to sync the Merkle root
	syncMerkleJob := NewSyncMerkleJob(aj.ctx, aj.repo, aj.smartContract)
	if err := aj.jobManager.AddJob("syncMerkleRoot", "@every 10s", syncMerkleJob); err != nil {
		log.Printf("Failed to add syncMerkleRoot job: %v", err)
	}

	aj.jobManager.Start()

	// Log running jobs
	log.Printf("Running jobs: %v", aj.jobManager.GetRunningJobs())

}

func (aj *AsyncJob) GetRunningJobs() []string {
	return aj.jobManager.GetRunningJobs()
}

func (aj *AsyncJob) Stop() {
	aj.jobManager.Stop()
}
