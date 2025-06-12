package ociloader

import (
	"context"
	"strings"
	"time"

	"github.com/konflux-ci/quality-dashboard/pkg/storage/ent/db"
)

// runWorkers initializes a fixed number of worker goroutines and distributes the given jobs among them.
// It blocks until all workers have completed processing.
func (l *Loader) runWorkers(ctx context.Context, jobs map[string]*repoJob) error {
	jobChan := make(chan *repoJob, len(jobs))
	done := make(chan struct{})

	// Start worker goroutines
	for i := 0; i < l.cfg.NumWorkers; i++ {
		go l.worker(ctx, i+1, jobChan, done)
	}

	// Dispatch all jobs to the channel
	for _, job := range jobs {
		jobChan <- job
	}
	close(jobChan)

	// Wait for all workers to complete
	for i := 0; i < l.cfg.NumWorkers; i++ {
		<-done
	}

	return nil
}

// worker processes a stream of repoJob items from the job channel.
// Each job represents a repository and its associated OCI artifacts.
// The worker only processes repositories whose latest artifact is considered stale
// based on the ProcessingThreshold defined in the configuration.
func (l *Loader) worker(ctx context.Context, id int, jobChan <-chan *repoJob, done chan<- struct{}) {
	for job := range jobChan {
		var latest time.Time
		var urls []string

		// Determine the most recent update timestamp among the artifacts
		for _, a := range job.Artifacts {
			urls = append(urls, strings.TrimPrefix(a.ArtifactURL, "quay.io/"))
			if a.UpdatedAt != nil && a.UpdatedAt.After(latest) {
				latest = *a.UpdatedAt
			} else if a.CreatedAt.After(latest) {
				latest = a.CreatedAt
			}
		}

		// Skip processing if the most recent update is within the freshness threshold
		if time.Since(latest) < l.cfg.ProcessingThreshold {
			l.logger.Sugar().Infof("[Worker %d] Skipping recent repository '%s'", id, job.Repository.RepositoryName)
			continue
		}

		l.logger.Sugar().Infof("[Worker %d] Processing '%s'", id, job.Repository.RepositoryName)

		// Attempt to process the repository with the OCI controller
		if errs := l.ociCtl.ProcessRepositories(urls, time.Since(latest)); len(errs) > 0 {
			l.logger.Sugar().Errorf("[Worker %d] OCI errors: %v", id, errs)
			continue
		}

		// Mark artifacts as updated in the storage layer
		for _, artifact := range job.Artifacts {
			if _, err := l.storage.UpdateOCIArtifact(ctx, artifact); err != nil {
				l.logger.Sugar().Errorf("[Worker %d] Failed to update artifact %s: %v", id, artifact.ID, err)
			}
		}
	}
	done <- struct{}{}
}

// groupByRepository takes a slice of OCI artifacts and groups them into repoJob structs,
// where each group corresponds to a unique repository.
func groupByRepository(artifacts []*db.OCI) map[string]*repoJob {
	jobMap := make(map[string]*repoJob)

	for _, artifact := range artifacts {
		if artifact.Edges.Oci == nil {
			continue // Skip orphaned artifacts with no repository linkage
		}
		repo := artifact.Edges.Oci
		if _, ok := jobMap[repo.ID]; !ok {
			jobMap[repo.ID] = &repoJob{Repository: repo}
		}
		jobMap[repo.ID].Artifacts = append(jobMap[repo.ID].Artifacts, artifact)
	}

	return jobMap
}
