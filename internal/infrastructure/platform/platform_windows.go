//go:build windows

package platform

import (
	"errors"
	"os/exec"
)

func (Services) MoveToTrash(paths []string) error {
	return errors.New("move to trash requires Windows shell integration; permanent delete is intentionally disabled")
}

func (Services) RevealInFileManager(path string) error {
	return exec.Command("explorer.exe", "/select,", path).Start()
}
func (Services) OpenFile(path string) error {
	return exec.Command("rundll32.exe", "url.dll,FileProtocolHandler", path).Start()
}
