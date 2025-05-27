package macro

import (
	"os"
	"path/filepath"

	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func ComposerDumpEnv() {
	p := project.GetProject()

	if p.HasDependency("symfony/flex") {
		logger := shell.GetLogger()
		pd := p.GetDirectory()
		envFile := filepath.Join(pd, ".env.local.php")

		if _, err := os.Stat(envFile); os.IsNotExist(err) {
			cmd := shell.NewComposerCommand("dump-env", "--no-interaction", "--no-scripts")

			if err := cmd.WithLogger(logger).Run(); err != nil {
				panic(errors.Wrap(err, "Failed to dump environment"))
			}
		}

	}
}
