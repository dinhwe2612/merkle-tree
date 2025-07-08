package cronjob

import (
	"context"
	"log"
	"merkle_module/domain/repo"
)

type AsyncJob struct {
	ctx        context.Context
	jobManager *JobManager
	repo       repo.Merkle
}

func NewAsyncJob(ctx context.Context, repo repo.Merkle) *AsyncJob {
	jobManager := NewJobManager()
	return &AsyncJob{
		ctx:        ctx,
		jobManager: jobManager,
		repo:       repo,
	}
}

func (aj *AsyncJob) Start() {
	// Add a job to sync the Merkle root
	syncMerkleJob := NewSyncMerkleJob(aj.ctx, aj.repo)
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
