package macro

import (
	"path/filepath"

	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func CheckRequirements() {
	logger := shell.GetLogger()
	p := project.GetProject()
	env := p.GetEnvironment()

	if env.IsProd() {
		if p.HasDependency("symfony/requirements-checker") {
			cmd := shell.NewPHPCommand(filepath.Join(p.GetDirectory(), "vendor", "bin", "requirements-checker"))

			if err := cmd.WithLogger(logger).WithLogLevel(0).Run(); err != nil {
				logger.Errorln(" [ERROR]")
				panic(errors.Wrap(err, "failed to execute requirements-checker command"))
			}
		}
	} else if env.IsDev() {
		cmd, err := shell.NewSymfonyCommand("check:requirements")

		if err != nil {
			panic(errors.Wrap(err, "failed to create symfony check:requirements command"))
		}

		if err := cmd.WithLogger(logger).WithLogLevel(0).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute symfony check:requirements command"))
		}
	}
}
