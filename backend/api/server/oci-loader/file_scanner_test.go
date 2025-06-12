package ociloader

import (
	"os"
	"path/filepath"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

// InitTestLogger creates a zap.Logger suitable for test output.
func InitTestLogger(t *testing.T) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel) // or InfoLevel
	logger, err := cfg.Build()
	if err != nil {
		t.Fatalf("failed to initialize zap logger: %v", err)
	}
	return logger
}

func TestScanForArtifacts_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	logger := InitTestLogger(t)

	sets, err := ScanForArtifacts(tmpDir, logger)

	assert.NoError(t, err)
	assert.Len(t, sets, 0)
}

func TestScanForArtifacts_MissingDirectory(t *testing.T) {
	logger := InitTestLogger(t)

	sets, err := ScanForArtifacts("/non/existent/path", logger)

	assert.NoError(t, err)
	assert.Nil(t, sets)
}

func TestScanForArtifacts_ValidArtifacts(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "run1")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "pipeline-status.json"), []byte("{}"), 0644)
	os.WriteFile(filepath.Join(subDir, "e2e-report.xml"), []byte("<testsuites></testsuites>"), 0644)

	logger := InitTestLogger(t)

	sets, err := ScanForArtifacts(tmpDir, logger)

	assert.NoError(t, err)
	assert.Len(t, sets, 1)
	assert.Equal(t, filepath.Join(subDir, "pipeline-status.json"), sets[0].PipelineStatusPath)
	assert.Equal(t, filepath.Join(subDir, "e2e-report.xml"), sets[0].E2EReportPath)
}

func TestScanForArtifacts_OnlyPipelineStatus(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "run2")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "pipeline-status.json"), []byte("{}"), 0644)

	logger := InitTestLogger(t)

	sets, err := ScanForArtifacts(tmpDir, logger)

	assert.NoError(t, err)
	assert.Len(t, sets, 1)
	assert.Equal(t, filepath.Join(subDir, "pipeline-status.json"), sets[0].PipelineStatusPath)
	assert.Equal(t, "", sets[0].E2EReportPath)
}
