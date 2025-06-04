package env

import (
	"path/filepath"

	"github.com/SoureCode/kyx/internal/php"
	"github.com/pkg/errors"
)

func loadPHPFile(envMap map[string]string, directory string, filename string) error {
	filePath := filepath.Join(directory, filename)
	parsedEnvMap, err := php.LoadFileDump(filePath)

	if err != nil {
		return errors.Wrapf(err, "failed to load PHP file: %s", filePath)
	}

	mergeEnvMaps(envMap, parsedEnvMap)

	return nil
}
