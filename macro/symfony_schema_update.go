package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonySchemaUpdate() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("doctrine/doctrine-bundle") {
		cmd := shell.NewConsoleCommand("doctrine:schema:update", "--dump-sql", "--no-interaction")

		if err := cmd.WithLogger(logger).WithLogLevel(3).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to dump schema"))
		}

		cmd = shell.NewConsoleCommand("doctrine:schema:update", "--force", "--no-interaction")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to update schema"))
		}
	}
}
