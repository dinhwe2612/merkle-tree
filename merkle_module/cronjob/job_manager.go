package cronjob

import (
	"log"
	"sync"

	"github.com/robfig/cron/v3"
)

// Job defines the interface for all cron jobs
type Job interface {
	Run()
}

// JobManager manages cron jobs
type JobManager struct {
	cron      *cron.Cron
	jobs      map[string]cron.EntryID
	mutex     sync.RWMutex
	isRunning bool
}

// NewJobManager creates a new JobManager
func NewJobManager() *JobManager {
	return &JobManager{
		cron:      cron.New(),
		jobs:      make(map[string]cron.EntryID),
		isRunning: false,
	}
}

// Start begins the cron scheduler
func (jm *JobManager) Start() {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()
	if !jm.isRunning {
		jm.cron.Start()
		jm.isRunning = true
		log.Println("JobManager started")
	}
}

// Stop halts the cron scheduler
func (jm *JobManager) Stop() {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()
	if jm.isRunning {
		jm.cron.Stop()
		jm.isRunning = false
		log.Println("JobManager stopped")
	}
}

// AddJob adds a new cron job with a specified schedule
func (jm *JobManager) AddJob(name, schedule string, job Job) error {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()

	entryID, err := jm.cron.AddFunc(schedule, job.Run)
	if err != nil {
		return err
	}
	jm.jobs[name] = entryID
	log.Printf("Added job: %s with schedule: %s", name, schedule)
	return nil
}

// RemoveJob removes a cron job
func (jm *JobManager) RemoveJob(name string) {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()

	if entryID, exists := jm.jobs[name]; exists {
		jm.cron.Remove(entryID)
		delete(jm.jobs, name)
		log.Printf("Removed job: %s", name)
	}
}

// GetRunningJobs returns a list of names of all currently running jobs
func (jm *JobManager) GetRunningJobs() []string {
	jm.mutex.RLock()
	defer jm.mutex.RUnlock()

	names := make([]string, 0, len(jm.jobs))
	for name := range jm.jobs {
		names = append(names, name)
	}
	return names
}
