package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyWorkerStop() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("symfony/messenger") {
		cmd := shell.NewConsoleCommand("messenger:stop-workers", "--no-interaction")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to stop Symfony Messenger workers"))
		}
	}
}
