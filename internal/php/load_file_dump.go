package php

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/pkg/errors"
)

func LoadFileDump(f string) (map[string]string, error) {
	info, err := os.Stat(f)

	if err != nil {
		if os.IsNotExist(err) {
			return map[string]string{}, nil
		}

		return nil, errors.Wrap(err, "failed to stat file")
	}

	if info.IsDir() {
		return nil, errors.Errorf("expected file but found directory")
	}

	dir := filepath.Dir(f)

	var buf bytes.Buffer

	script := fmt.Sprintf(`$entries = include %q; if (!is_array($entries)) exit(1); foreach ($entries as $key => $value) { echo "$key=$value\n"; }`, f)
	cmd := exec.Command("php", "-r", script)
	cmd.Dir = dir
	cmd.Stdout = &buf
	cmd.Stderr = &buf

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "failed to start php command")
	}

	if err := cmd.Wait(); err != nil {
		return nil, errors.Wrapf(err, "PHP script execution failed: %s", strings.TrimSpace(buf.String()))
	}

	parsedEnvMap, err := godotenv.UnmarshalBytes(buf.Bytes())

	if err != nil {
		return nil, errors.Wrap(err, "failed to parse environment variables")
	}

	return parsedEnvMap, nil
}
