package commands

import (
	"github.com/SoureCode/kyx/macro"
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
	"github.com/symfony-cli/console"
)

var deploymentCommand = &console.Command{
	Name:        "deployment",
	Description: "Deploy the application",
	Action: func(c *console.Context) error {
		p := project.GetProject()
		env := p.GetEnvironment()

		if !env.IsProd() {
			return errors.New("Deployment can only be run in production environment")
		}

		macro.ComposerInstall()
		macro.CheckRequirements()
		macro.ComposerDumpEnv()

		// stop
		macro.SymfonyWorkerStop()
		macro.SoureCodeScreenStop()

		// database
		macro.WaitForDatabase()
		macro.SymfonyMigrationsMigrate()

		// cache and assets
		macro.SymfonyCacheClear()
		macro.SymfonyAssetsInstall()
		macro.SymfonyImportMapInstall()

		// start
		macro.SoureCodeScreenStart()

		shell.GetLogger().LogDuration()

		return nil
	},
}
