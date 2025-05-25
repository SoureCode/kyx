package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SymfonyCacheClear() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("symfony/framework-bundle") {
		cmd := shell.NewConsoleCommand("cache:clear", "--no-interaction", "--no-warmup")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to clear cache"))
		}

		cmd = shell.NewConsoleCommand("cache:warmup", "--no-interaction")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to warm up cache"))
		}
	}
}
