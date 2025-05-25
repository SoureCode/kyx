package macro

import (
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func SoureCodeScreenStart() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("sourecode/screen-bundle") {
		cmd := shell.NewConsoleCommand("screen:start", "--no-interaction")

		if err := cmd.WithLogger(logger).Run(); err != nil {
			panic(errors.Wrap(err, "failed to execute command to start screens"))
		}
	}
}
