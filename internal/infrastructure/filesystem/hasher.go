package filesystem

import (
	"context"
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
	return HashFileContext(context.Background(), path)
}

func HashFileContext(ctx context.Context, path string) (HashResult, error) {
	if err := ctx.Err(); err != nil {
		return HashResult{}, err
	}
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
	if _, err := io.CopyBuffer(h, contextReader{ctx: ctx, reader: f}, buf); err != nil {
		return HashResult{}, err
	}
	after, err := os.Stat(path)
	if err != nil {
		return HashResult{}, err
	}
	return HashResult{
		SHA256:   hex.EncodeToString(h.Sum(nil)),
		Unstable: !os.SameFile(before, after) || before.Size() != after.Size() || !before.ModTime().Equal(after.ModTime()),
	}, nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r contextReader) Read(buffer []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	return r.reader.Read(buffer)
}
