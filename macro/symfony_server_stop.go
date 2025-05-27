package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyServerStop() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("symfony/framework-bundle") {
		cmd, err := shell.NewSymfonyCommand("server:stop")

		if err != nil {
			panic(errors.Wrap(err, "failed to create Symfony command to stop server"))
		}

		if err = cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to stop server"))
		}
	}
}
