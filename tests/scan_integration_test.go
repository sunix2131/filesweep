package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"filesweep/internal/domain/settings"
	"filesweep/internal/infrastructure/filesystem"

	"github.com/stretchr/testify/require"
)

func TestScanFindsOnlyTrueDuplicates(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(dir, "nested"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a.pdf"), []byte("same"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "a-copy.pdf"), []byte("same"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "same-size.txt"), []byte("diff"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "nested", "same-size-2.txt"), []byte("xxxx"), 0o644))
	result, err := (filesystem.Scanner{}).Scan(context.Background(), []string{dir}, settings.Default())
	require.NoError(t, err)
	candidates := filesystem.DuplicateCandidates(result.Files)
	candidates = filesystem.ApplyHashes(candidates, 2, filesystem.HashFile)
	groups := filesystem.BuildDuplicateGroups(result.Session.ID, candidates)
	require.Len(t, groups, 1)
	require.Equal(t, int64(4), groups[0].EstimatedReclaimableBytes)
}
