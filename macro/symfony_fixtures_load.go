package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyFixturesLoad() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDevDependency("doctrine/doctrine-fixtures-bundle") {
		cmd := shell.NewConsoleCommand("doctrine:fixtures:load", "--no-interaction")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to load fixtures"))
		}
	}
}
