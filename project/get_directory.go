package project

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"

	"github.com/SoureCode/kyx/git"
)

func getDirectory() (string, error) {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "could not get current working directory")
	}

	// Try to find git root
	gitRoot, err := git.RootDirectory(currentWorkingDirectory)

	if err == nil && gitRoot != "" {
		composerPath := filepath.Join(gitRoot, "composer.json")

		if _, err := os.Stat(composerPath); err == nil {
			return gitRoot, nil
		}

		return "", errors.Wrap(err, fmt.Sprintf("git root found at '%s' but no composer.json found there", gitRoot))
	}

	// Fallback to upward search
	if _, err := os.Stat(filepath.Join(currentWorkingDirectory, "composer.json")); err == nil {
		return currentWorkingDirectory, nil
	}

	currentDirectory := currentWorkingDirectory

	for {
		if _, err := os.Stat(filepath.Join(currentDirectory, "composer.json")); err == nil {
			return currentDirectory, nil
		}

		parent := filepath.Dir(currentDirectory)
		if parent == currentDirectory {
			return "", errors.Wrap(err, fmt.Sprintf("could not find project directory: searched from '%s' upwards for 'composer.json'", currentWorkingDirectory))
		}

		currentDirectory = parent
	}
}
