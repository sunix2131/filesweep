//go:build darwin

package platform

import (
	"os/exec"
	"strconv"
)

func (Services) MoveToTrash(paths []string) error {
	args := append([]string{"-e"}, trashScript(paths...))
	return exec.Command("osascript", args...).Run()
}

func (Services) RevealInFileManager(path string) error { return exec.Command("open", "-R", path).Run() }
func (Services) OpenFile(path string) error            { return exec.Command("open", path).Start() }

func trashScript(paths ...string) string {
	s := `tell application "Finder"` + "\n"
	for _, p := range paths {
		s += `delete POSIX file ` + strconv.Quote(p) + "\n"
	}
	return s + "end tell"
}
