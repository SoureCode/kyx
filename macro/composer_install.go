package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func ComposerInstall(args ...string) {
	p := project.GetProject()
	env := p.GetEnvironment()

	defaultArgs := []string{"install", "--no-interaction", "--no-scripts"}

	if env.IsProd() {
		defaultArgs = append(defaultArgs, "--no-dev", "--optimize-autoloader")
	}

	args = append(defaultArgs, args...)

	cmd := shell.NewComposerCommand(args...)

	if err := cmd.Run(); err != nil {
		panic(errors.Wrap(err, "failed to execute composer install command"))
	}
}
