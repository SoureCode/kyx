package tools

import (
	"os"
	"path/filepath"

	"github.com/SoureCode/kyx/project"
)

func GetDirectory() string {
	p := project.GetProject()
	toolsDirectory := filepath.Join(p.GetDirectory(), "tools")

	if _, err := os.Stat(toolsDirectory); os.IsNotExist(err) {
		if err := os.Mkdir(toolsDirectory, 0755); err != nil {
			panic(err)
		}
	}

	return toolsDirectory
}
