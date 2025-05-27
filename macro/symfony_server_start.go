package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyServerStart() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("symfony/framework-bundle") {
		cmd, err := shell.NewSymfonyCommand("serve", "--daemon")

		if err != nil {
			panic(errors.Wrap(err, "failed to create Symfony command to start server"))
		}

		if err = cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to start server"))
		}
	}
}
