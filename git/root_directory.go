package git

import (
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

// RootDirectory tries to find the git root directory from the given path.
func RootDirectory(startDir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = startDir
	out, err := cmd.Output()

	if err != nil {
		return "", errors.Wrapf(err, "failed to read git root directory")
	}

	return string(filepath.Clean(string(out[:len(out)-1]))), nil
}
