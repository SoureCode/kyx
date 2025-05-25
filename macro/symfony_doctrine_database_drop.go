package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyDoctrineDatabaseDrop() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("doctrine/doctrine-bundle") {
		cmd := shell.NewConsoleCommand("doctrine:database:drop", "--no-interaction", "--if-exists", "--force")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to drop database"))
		}

		cmd = shell.NewConsoleCommand("doctrine:database:create", "--no-interaction", "--if-not-exists")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to create database"))
		}
	}
}
