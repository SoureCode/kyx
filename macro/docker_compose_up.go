package macro

import (
	"os"
	"path/filepath"

	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func anyFileExists(filenames ...string) bool {
	for _, filename := range filenames {
		if _, err := os.Stat(filename); err == nil {
			return true
		}
	}

	return false
}

func DockerComposeUp() {
	logger := shell.GetLogger()
	p := project.GetProject()
	pd := p.GetDirectory()
	files := []string{
		filepath.Join(pd, "compose.yaml"),
		filepath.Join(pd, "compose.yml"),
		filepath.Join(pd, "docker-compose.yaml"),
		filepath.Join(pd, "docker-compose.yml"),
	}

	if anyFileExists(files...) {
		cmd := shell.NewDockerCommand("compose", "up", "--detach", "--remove-orphans")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute docker compose up command"))
		}
	}
}
