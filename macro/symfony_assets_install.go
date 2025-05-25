package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyAssetsInstall() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("symfony/framework-bundle") {
		cmd := shell.NewConsoleCommand("assets:install", "--no-interaction")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to install assets"))
		}
	}
}
