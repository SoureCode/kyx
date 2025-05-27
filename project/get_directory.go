package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SoureCode/kyx/git"
	"github.com/pkg/errors"
)

func getDirectory() string {
	currentWorkingDirectory, err := os.Getwd()
	if err != nil {
		panic(errors.Wrap(err, "could not get current working directory"))
	}

	// Try to find git root
	gitRoot, err := git.RootDirectory(currentWorkingDirectory)

	if err == nil && gitRoot != "" {
		composerPath := filepath.Join(gitRoot, "composer.json")

		if _, err := os.Stat(composerPath); err == nil {
			return gitRoot
		}

		panic(fmt.Sprintf("git root found at '%s' but no composer.json found there", gitRoot))
	}

	// Fallback to upward search
	if _, err := os.Stat(filepath.Join(currentWorkingDirectory, "composer.json")); err == nil {
		return currentWorkingDirectory
	}

	currentDirectory := currentWorkingDirectory

	for {
		if _, err := os.Stat(filepath.Join(currentDirectory, "composer.json")); err == nil {
			return currentDirectory
		}

		parent := filepath.Dir(currentDirectory)
		if parent == currentDirectory {
			panic(fmt.Sprintf("could not find project directory: searched from '%s' upwards for 'composer.json'", currentWorkingDirectory))
		}

		currentDirectory = parent
	}
}
