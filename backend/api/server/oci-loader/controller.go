package ociloader

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/konflux-ci/qe-tools/pkg/oci"
	"github.com/konflux-ci/quality-dashboard/pkg/logger"
	"github.com/konflux-ci/quality-dashboard/pkg/storage"
	"go.uber.org/zap"
)

// Loader is responsible for managing the full lifecycle of OCI artifact processing:
// - Downloading and caching OCI content
// - Grouping repositories and invoking workers
// - Parsing results and storing test outcomes
// - Cleaning up execution directories (if enabled)
type Loader struct {
	// Configuration for loader execution (concurrency, paths, thresholds)
	cfg LoaderConfig

	// Structured logger for debug and audit logs
	logger *zap.Logger

	// Interface to the backing storage layer
	storage storage.Storage

	// Unique identifier for this run
	execID uuid.UUID

	// OCI artifact processor and extractor
	ociCtl *oci.Controller

	// Temporary directory to store downloaded artifacts
	artifactDir string

	// Temporary cache directory for OCI blovs
	cacheDir string
}

// NewLoader constructs a new Loader with a fresh execution ID,
// prepares isolated artifact/cache directories, and initializes the download of oci artifacts.
// NOTE! THe OCI artifacts only used if is Konflux CI
func NewLoader(cfg LoaderConfig, storage storage.Storage) *Loader {
	execID := uuid.New()
	logger, _ := logger.InitZap("info")

	artifactDir := filepath.Join(cfg.ArtifactsDir, execID.String())
	cacheDir := filepath.Join(cfg.CacheDir, execID.String())

	return &Loader{
		cfg:         cfg,
		logger:      logger,
		storage:     storage,
		execID:      execID,
		artifactDir: artifactDir,
		cacheDir:    cacheDir,
	}
}

// Run executes the OCI artifact processing lifecycle in the following phases:
//  1. Optionally sets up deferred cleanup for working directories
//  2. Fetches artifact metadata from persistent storage
//  3. Groups artifacts by repository and processes them using a worker pool
//  4. Scans extracted artifacts for result files (pipeline status and xUnit test data)
//  5. Parses and persists results back into the storage layer
//
// It returns an error if any critical step in the lifecycle fails.
func (l *Loader) Run(ctx context.Context) error {
	if l.cfg.CleanupOnCompletion {
		defer func() {
			if err := os.RemoveAll(l.artifactDir); err != nil {
				l.logger.Sugar().Errorf("Failed to remove artifact directory %s: %v", l.artifactDir, err)
			}
			if err := os.RemoveAll(l.cacheDir); err != nil {
				l.logger.Sugar().Errorf("Failed to remove cache directory %s: %v", l.cacheDir, err)
			}
		}()
	}

	artifacts, err := l.storage.GetAllOCIArtifacts()
	if err != nil {
		return fmt.Errorf("Failed to get OCI artifacts: %w", err)
	}
	l.logger.Sugar().Infof("Found %d artifacts", len(artifacts))

	jobs := groupByRepository(artifacts)

	l.ociCtl, err = oci.NewController(l.artifactDir, l.cacheDir)
	if err != nil {
		return fmt.Errorf("Failed to create OCI controller: %w", err)
	}

	if err := l.runWorkers(ctx, jobs); err != nil {
		return err
	}

	sets, err := ScanForArtifacts(l.artifactDir, l.logger)
	if err != nil {
		return fmt.Errorf("Failed to scan artifact sets: %w", err)
	}

	for _, set := range sets {
		if err := ProcessArtifactSet(set, l.storage, l.logger); err != nil {
			l.logger.Sugar().Errorf("Failed to process artifact set %s: %v", set.Directory, err)
		}
	}

	return nil
}
