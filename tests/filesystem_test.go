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
