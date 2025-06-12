package ociloader

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

// ScanForArtifacts recursively scans the given root directory for directories containing
// Tekton pipeline result files. A valid artifact set is defined by the presence of a
// 'pipeline-status.json' file, and optionally an 'e2e-report.xml' file.
//
// The function returns a slice of ArtifactSet structs, each representing a directory where
// a complete or partial artifact set was found.
func ScanForArtifacts(root string, logger *zap.Logger) ([]ArtifactSet, error) {
	var sets []ArtifactSet
	dirs := map[string]map[string]bool{}

	if _, err := os.Stat(root); os.IsNotExist(err) {
		logger.Sugar().Infof("Directory '%s' doesn't exist", root)
		return nil, nil
	}

	// TODO: Make the target files customizable by environments or flags
	targets := map[string]bool{
		"pipeline-status.json": true,
		"e2e-report.xml":       true,
	}

	// Walk the directory tree and track which files exist in each subdirectory
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		if !targets[d.Name()] {
			return nil
		}
		dir := filepath.Dir(path)
		if dirs[dir] == nil {
			dirs[dir] = map[string]bool{}
		}
		dirs[dir][d.Name()] = true
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk error: %w", err)
	}

	for dir, files := range dirs {
		if files["pipeline-status.json"] {
			set := ArtifactSet{
				Directory:          dir,
				PipelineStatusPath: filepath.Join(dir, "pipeline-status.json"),
			}
			if files["e2e-report.xml"] {
				set.E2EReportPath = filepath.Join(dir, "e2e-report.xml")
			}
			sets = append(sets, set)
		}
	}

	if len(sets) == 0 {
		logger.Warn("No valid artifact sets found")
	}

	return sets, nil
}
