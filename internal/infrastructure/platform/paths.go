package platform

import (
	"os"
	"path/filepath"
)

func (Services) GetAppDataDir() (string, error) {
	root, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "FileSweep"), nil
}

func (Services) GetCacheDir() (string, error) {
	root, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "FileSweep"), nil
}

func (Services) GetLogsDir() (string, error) {
	dir, err := Services{}.GetAppDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "logs"), nil
}
