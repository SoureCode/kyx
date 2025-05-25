package commands

import (
	"github.com/SoureCode/kyx/macro"
	"github.com/SoureCode/kyx/shell"
	"github.com/symfony-cli/console"
)

var stopCommand = &console.Command{
	Name:        "stop",
	Description: "Stop the project",
	Action: func(c *console.Context) error {
		macro.SymfonyWorkerStop()
		macro.SoureCodeScreenStop()
		macro.DockerComposeDown()
		macro.SymfonyServerStop()

		shell.GetLogger().LogDuration()

		return nil
	},
}
