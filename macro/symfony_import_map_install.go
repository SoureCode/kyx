package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyImportMapInstall() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("symfony/asset-mapper") {
		cmd := shell.NewConsoleCommand("importmap:install", "--no-interaction")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to install import map"))
		}

		env := p.GetEnvironment()

		if env.IsProd() {
			cmd = shell.NewConsoleCommand("asset-map:compile", "--no-interaction")

			if err := cmd.WithLogger(logger).Run(); err != nil {
				panic(errors.Wrap(err, "failed to execute command to compile asset map"))
			}
		}
	}
}
