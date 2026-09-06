package filesystem

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func AutoRenamePath(target string) string {
	if _, err := os.Stat(target); errors.Is(err, os.ErrNotExist) {
		return target
	}
	ext := filepath.Ext(target)
	base := strings.TrimSuffix(target, ext)
	for i := 1; i < 10000; i++ {
		candidate := fmt.Sprintf("%s (%d)%s", base, i, ext)
		if _, err := os.Stat(candidate); errors.Is(err, os.ErrNotExist) {
			return candidate
		}
	}
	return target
}

func SafeMove(source, target string, expectedSize int64, expectedSHA string) error {
	source = filepath.Clean(source)
	target = filepath.Clean(target)
	if source == target {
		return errors.New("source and target paths are the same")
	}

	info, err := VerifyScannedFile(source, expectedSize, expectedSHA)
	if err != nil {
		return err
	}
	target = AutoRenamePath(target)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	// A hard link is an atomic no-overwrite move on the same filesystem.
	// If it is not supported, the copy path below handles another volume.
	if err := os.Link(source, target); err == nil {
		if err := os.Remove(source); err != nil {
			_ = os.Remove(target)
			return fmt.Errorf("remove source after linking: %w", err)
		}
		return nil
	}

	tmp, err := copyToTemporary(source, filepath.Dir(target), info.Mode().Perm())
	if err != nil {
		return err
	}
	defer os.Remove(tmp)

	copied, err := os.Stat(tmp)
	if err != nil {
		return err
	}
	if copied.Size() != info.Size() {
		return errors.New("copied file size mismatch")
	}
	if expectedSHA != "" {
		res, err := HashFile(tmp)
		if err != nil {
			return err
		}
		if res.Unstable || res.SHA256 != expectedSHA {
			return errors.New("copied file hash mismatch")
		}
	}
	if err := os.Chtimes(tmp, info.ModTime(), info.ModTime()); err != nil {
		return fmt.Errorf("preserve modification time: %w", err)
	}
	if err := os.Link(tmp, target); err != nil {
		return fmt.Errorf("create destination: %w", err)
	}
	if err := os.Remove(source); err != nil {
		_ = os.Remove(target)
		return fmt.Errorf("remove source after copying: %w", err)
	}
	return nil
}

// VerifyScannedFile prevents an action from operating on a path whose file was
// replaced or modified after the scan result was created.
func VerifyScannedFile(path string, expectedSize int64, expectedSHA string) (os.FileInfo, error) {
	info, err := os.Lstat(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, errors.New("source is not a regular file")
	}
	if expectedSize >= 0 && info.Size() != expectedSize {
		return nil, errors.New("file changed after scan")
	}
	if expectedSHA != "" {
		hash, err := HashFile(path)
		if err != nil {
			return nil, err
		}
		if hash.Unstable || hash.SHA256 != expectedSHA {
			return nil, errors.New("file changed after scan")
		}
	}
	return info, nil
}

func UndoMove(source, target string, expectedSize int64, expectedSHA string) error {
	if _, err := os.Stat(source); err == nil {
		return errors.New("original path is not free")
	}
	info, err := os.Stat(target)
	if err != nil {
		return err
	}
	if expectedSize >= 0 && info.Size() != expectedSize {
		return errors.New("moved file changed after operation")
	}
	if expectedSHA != "" {
		res, err := HashFile(target)
		if err != nil {
			return err
		}
		if res.SHA256 != expectedSHA {
			return errors.New("moved file hash changed after operation")
		}
	}
	return SafeMove(target, source, expectedSize, expectedSHA)
}

func copyToTemporary(source, targetDir string, mode os.FileMode) (string, error) {
	in, err := os.Open(source)
	if err != nil {
		return "", err
	}
	defer in.Close()
	out, err := os.CreateTemp(targetDir, ".filesweep-copying-*")
	if err != nil {
		return "", err
	}
	path := out.Name()
	remove := true
	defer func() {
		_ = out.Close()
		if remove {
			_ = os.Remove(path)
		}
	}()
	if err := out.Chmod(mode); err != nil {
		return "", err
	}
	if _, err := io.CopyBuffer(out, in, make([]byte, 2*1024*1024)); err != nil {
		return "", err
	}
	if err := out.Sync(); err != nil {
		return "", err
	}
	if err := out.Close(); err != nil {
		return "", err
	}
	remove = false
	return path, nil
}
