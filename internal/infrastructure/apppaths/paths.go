package apppaths

import (
	"os"
	"path/filepath"
)

type Paths struct {
	ConfigDir     string
	CacheDir      string
	LogsDir       string
	ExportsDir    string
	ThumbnailsDir string
	BackupsDir    string
	DBPath        string
}

func Resolve() (Paths, error) {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		return Paths{}, err
	}
	cacheRoot, err := os.UserCacheDir()
	if err != nil {
		return Paths{}, err
	}
	p := Paths{
		ConfigDir: filepath.Join(configRoot, "FileSweep"),
		CacheDir:  filepath.Join(cacheRoot, "FileSweep"),
	}
	p.LogsDir = filepath.Join(p.ConfigDir, "logs")
	p.ExportsDir = filepath.Join(p.ConfigDir, "exports")
	p.ThumbnailsDir = filepath.Join(p.CacheDir, "thumbnails")
	p.BackupsDir = filepath.Join(p.ConfigDir, "backups")
	p.DBPath = filepath.Join(p.ConfigDir, "filesweep.db")
	for _, dir := range []string{p.ConfigDir, p.CacheDir, p.LogsDir, p.ExportsDir, p.ThumbnailsDir, p.BackupsDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return Paths{}, err
		}
	}
	return p, nil
}
