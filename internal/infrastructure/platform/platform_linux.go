//go:build linux

package platform

import (
	"errors"
	"os/exec"
)

func (Services) MoveToTrash(paths []string) error {
	args := append([]string{"trash"}, paths...)
	if err := exec.Command("gio", args...).Run(); err == nil {
		return nil
	}
	return errors.New("system trash is unavailable")
}

func (Services) RevealInFileManager(path string) error { return exec.Command("xdg-open", path).Start() }
func (Services) OpenFile(path string) error            { return exec.Command("xdg-open", path).Start() }
