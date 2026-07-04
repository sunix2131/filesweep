package filesystem

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type HashResult struct {
	SHA256   string
	Unstable bool
}

func HashFile(path string) (HashResult, error) {
	before, err := os.Stat(path)
	if err != nil {
		return HashResult{}, err
	}
	f, err := os.Open(path)
	if err != nil {
		return HashResult{}, err
	}
	defer f.Close()
	h := sha256.New()
	buf := make([]byte, 2*1024*1024)
	if _, err := io.CopyBuffer(h, f, buf); err != nil {
		return HashResult{}, err
	}
	after, err := os.Stat(path)
	if err != nil {
		return HashResult{}, err
	}
	return HashResult{
		SHA256:   hex.EncodeToString(h.Sum(nil)),
		Unstable: before.Size() != after.Size() || !before.ModTime().Equal(after.ModTime()),
	}, nil
}
