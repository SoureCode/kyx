package tools

import (
	"path/filepath"

	"github.com/SoureCode/kyx/shell"
)

func NewToolCommand(toolName string, mapping Mapping, args ...string) *shell.Command {
	directory := GetDirectory()
	install(toolName, mapping)
	toolBinary := filepath.Join(directory, toolName)

	return shell.NewPHPCommand(append([]string{toolBinary}, args...)...).WithPassthrough()
}
