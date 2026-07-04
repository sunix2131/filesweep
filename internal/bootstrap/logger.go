package bootstrap

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"filesweep/internal/infrastructure/apppaths"
)

func NewLogger(paths apppaths.Paths) (*slog.Logger, io.Closer, error) {
	name := "filesweep-" + time.Now().Format("2006-01-02") + ".log"
	f, err := os.OpenFile(filepath.Join(paths.LogsDir, name), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, err
	}
	return slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelInfo})), f, nil
}
