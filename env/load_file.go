package env

import (
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

func loadFile(envMap map[string]string, directory string, filename string) error {
	filePath := filepath.Join(directory, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // File does not exist, skip loading
	}

	readEnvMap, err := godotenv.Read(filePath)

	if err != nil {
		return errors.Wrapf(err, "failed to read environment file %s", filePath)
	}

	mergeEnvMaps(envMap, readEnvMap)

	return nil
}
