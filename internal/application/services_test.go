package application

import (
	"os"
	"path/filepath"
	"testing"

	"filesweep/internal/domain/actions"

	"github.com/stretchr/testify/require"
)

func TestValidateScanPathsRejectsFilesAndDeduplicatesFolders(t *testing.T) {
	directory := t.TempDir()
	file := filepath.Join(directory, "file.txt")
	require.NoError(t, os.WriteFile(file, []byte("test"), 0o600))

	paths, err := validateScanPaths([]string{directory, directory})
	require.NoError(t, err)
	require.Equal(t, []string{directory}, paths)

	_, err = validateScanPaths([]string{file})
	require.ErrorContains(t, err, "not a folder")
}

func TestValidateActionRejectsIncompletePlans(t *testing.T) {
	require.Error(t, validateAction(actions.TypeMoveToFolder, nil))
	require.Error(t, validateAction("delete", []actions.Item{{SourcePath: "/tmp/a", SourceSizeBytes: 1}}))
	require.Error(t, validateAction(actions.TypeMoveToFolder, []actions.Item{{SourcePath: "/tmp/a", SourceSizeBytes: 1}}))
	require.NoError(t, validateAction(actions.TypeMoveToTrash, []actions.Item{{SourcePath: "/tmp/a", SourceSizeBytes: 1}}))
}
