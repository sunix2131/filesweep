package tests

import (
	"os"
	"path/filepath"
	"testing"

	"filesweep/internal/infrastructure/filesystem"

	"github.com/stretchr/testify/require"
)

func TestCategoryAndHidden(t *testing.T) {
	category, mimeType := filesystem.CategoryFor("photo.jpg")
	require.Equal(t, "images", category)
	require.Contains(t, mimeType, "image/")
	require.True(t, filesystem.IsHidden(".secret"))
}

func TestHashFileDetectsStableDuplicate(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.pdf")
	b := filepath.Join(dir, "b.pdf")
	require.NoError(t, os.WriteFile(a, []byte("same"), 0o644))
	require.NoError(t, os.WriteFile(b, []byte("same"), 0o644))
	ha, err := filesystem.HashFile(a)
	require.NoError(t, err)
	hb, err := filesystem.HashFile(b)
	require.NoError(t, err)
	require.Equal(t, ha.SHA256, hb.SHA256)
	require.False(t, ha.Unstable)
}

func TestAutoRenamePath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.pdf")
	require.NoError(t, os.WriteFile(path, []byte("x"), 0o644))
	require.Equal(t, filepath.Join(dir, "file (1).pdf"), filesystem.AutoRenamePath(path))
}

func TestSafeMoveChecksHashBeforeMoving(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.txt")
	target := filepath.Join(dir, "target.txt")
	require.NoError(t, os.WriteFile(source, []byte("changed"), 0o640))

	err := filesystem.SafeMove(source, target, int64(len("changed")), "wrong-hash")

	require.ErrorContains(t, err, "changed after scan")
	require.FileExists(t, source)
	require.NoFileExists(t, target)
}

func TestSafeMoveDoesNotOverwriteExistingTarget(t *testing.T) {
	dir := t.TempDir()
	source := filepath.Join(dir, "source.txt")
	target := filepath.Join(dir, "target.txt")
	require.NoError(t, os.WriteFile(source, []byte("source"), 0o640))
	require.NoError(t, os.WriteFile(target, []byte("target"), 0o644))
	hash, err := filesystem.HashFile(source)
	require.NoError(t, err)

	require.ErrorContains(t, filesystem.SafeMove(source, target, int64(len("source")), hash.SHA256), "destination already exists")

	original, err := os.ReadFile(target)
	require.NoError(t, err)
	originalSource, err := os.ReadFile(source)
	require.NoError(t, err)
	require.Equal(t, []byte("target"), original)
	require.Equal(t, []byte("source"), originalSource)
	require.NoFileExists(t, filepath.Join(dir, "target (1).txt"))
}
