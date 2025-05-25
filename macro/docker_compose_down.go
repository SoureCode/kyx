package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
	"path/filepath"
)

func DockerComposeDown() {
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
		cmd := shell.NewDockerCommand("compose", "down", "--remove-orphans")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute docker compose down command"))
		}
	}
}
