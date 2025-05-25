package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyMigrationsMigrate() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("doctrine/doctrine-migrations-bundle") {
		cmd := shell.NewConsoleCommand("doctrine:migrations:migrate", "--no-interaction", "--allow-no-migration", "--all-or-nothing")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to migrate database"))
		}
	}
}
