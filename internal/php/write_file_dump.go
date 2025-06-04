package php

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

func WriteFileDump(file string, data map[string]string) error {
	content := formatData(data)
	err := os.WriteFile(file, []byte(content), 0644)
	return errors.Wrap(err, "failed to write data to file")
}

func formatData(data map[string]string) string {
	lines := []string{
		"<?php",
		"",
		"return [",
	}

	for key, value := range data {
		lines = append(lines, "    '"+key+"' => '"+value+"',")
	}

	lines = append(lines, "];")

	return strings.Join(lines, "\n") + "\n"
}
