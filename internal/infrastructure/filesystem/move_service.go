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
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if expectedSize >= 0 && info.Size() != expectedSize {
		return errors.New("file changed after scan")
	}
	target = AutoRenamePath(target)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}
	if err := os.Rename(source, target); err == nil {
		return nil
	}
	tmp := target + ".filesweep-copying"
	if err := copyFile(source, tmp); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	copied, err := os.Stat(tmp)
	if err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if copied.Size() != info.Size() {
		_ = os.Remove(tmp)
		return errors.New("copied file size mismatch")
	}
	if expectedSHA != "" {
		res, err := HashFile(tmp)
		if err != nil {
			_ = os.Remove(tmp)
			return err
		}
		if res.SHA256 != expectedSHA {
			_ = os.Remove(tmp)
			return errors.New("copied file hash mismatch")
		}
	}
	if err := os.Rename(tmp, target); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return os.Remove(source)
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

func copyFile(source, target string) error {
	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.CopyBuffer(out, in, make([]byte, 2*1024*1024)); err != nil {
		return err
	}
	return out.Sync()
}
