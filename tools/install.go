package tools

import (
	"github.com/pkg/errors"
	"github.com/symfony-cli/terminal"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func install(name string, mapping Mapping) string {
	logger := terminal.Logger
	url := mapping[name]
	toolsDirectory := GetDirectory()
	toolFilePath := filepath.Join(toolsDirectory, name)

	if _, err := os.Stat(toolFilePath); os.IsNotExist(err) {
		logger.Info().Msgf("Tool '%s' not found, downloading from %s", name, url)
		response, err := http.Get(url)

		if err != nil {
			logger.Error().Msgf("Error downloading tool '%s': %v", name, err)
			panic(errors.Wrap(err, "could not download tool"))
		}

		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.Error().Msgf("Error closing response body for tool '%s': %v", name, err)
				panic(errors.Wrap(err, "could not close response body"))
			}
		}(response.Body)

		if response.StatusCode != http.StatusOK {
			logger.Error().Msgf("Failed to download tool '%s': %s", name, response.Status)
			panic(errors.Errorf("failed to download tool: %s", response.Status))
		}

		logger.Info().Msgf("Creating file '%s' at %s", name, toolFilePath)
		file, err := os.Create(toolFilePath)
		if err != nil {
			logger.Error().Msgf("Error creating file '%s': %v", name, err)
			panic(errors.Wrap(err, "could not create tool file"))
		}

		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				logger.Error().Msgf("Error closing file '%s': %v", name, err)
				panic(errors.Wrap(err, "could not close file"))
			}
		}(file)

		logger.Info().Msgf("Writing to file '%s'", name)
		_, err = io.Copy(file, response.Body)
		if err != nil {
			logger.Error().Msgf("Error writing '%s' to file: %v", name, err)
			panic(errors.Wrap(err, "could not write file"))
		}

		logger.Info().Msgf("Setting executable permissions for file '%s'", name)
		err = os.Chmod(toolFilePath, 0755)

		if err != nil {
			logger.Error().Msgf("Error setting permissions for file '%s': %v", name, err)
			panic(errors.Wrap(err, "could not make file executable"))
		}
	}

	return toolFilePath
}
