package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

func loadPHPFile(envMap map[string]string, directory string, filename string) error {
	filePath := filepath.Join(directory, filename)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return errors.Wrapf(err, "failed to stat file %s", filePath)
	}

	if info.IsDir() {
		return errors.Errorf("expected file but found directory at %s", filePath)
	}

	script := fmt.Sprintf(`$entries = include %q; if (!is_array($entries)) exit(1); foreach ($entries as $key => $value) { echo "$key=$value\n"; }`, filePath)
	cmd := exec.Command("php", "-r", script)
	cmd.Dir = directory

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrapf(err, "failed to create stdout pipe for reading %s", filePath)
	}

	defer func(stdout io.ReadCloser) {
		err := stdout.Close()
		if err != nil {
			panic(errors.Wrapf(err, "failed to close stdout pipe for %s", filePath))
		}
	}(stdout)

	if err := cmd.Start(); err != nil {
		return errors.Wrapf(err, "failed to start command for reading %s", filePath)
	}

	parsedEnvMap, err := godotenv.Parse(stdout)

	if err != nil {
		return errors.Wrapf(err, "failed to parse environment variables from %s", filePath)
	}

	if err := cmd.Wait(); err != nil {
		return errors.Wrapf(err, "PHP script execution failed for %s", filePath)
	}

	mergeEnvMaps(envMap, parsedEnvMap)

	return nil
}
