package commands

import (
	"github.com/SoureCode/kyx/macro"
	"github.com/SoureCode/kyx/shell"
	"github.com/symfony-cli/console"
)

var startCommand = &console.Command{
	Name:        "start",
	Description: "Start the project",
	Flags: []console.Flag{
		&console.BoolFlag{
			Name:         "no-schema-update",
			DefaultValue: false,
			Usage:        "Skip schema update",
		},
		&console.BoolFlag{
			Name:         "no-fixtures-load",
			DefaultValue: false,
			Usage:        "Skip fixtures loading",
		},
	},
	Action: func(c *console.Context) error {
		macro.ComposerInstall()
		macro.CheckRequirements()

		macro.DockerComposeUp()
		macro.WaitForDatabase()

		macro.SymfonyMigrationsMigrate()

		if !c.Bool("no-schema-update") {
			macro.SymfonySchemaUpdate()
		}

		if !c.Bool("no-fixtures-load") {
			macro.SymfonyFixturesLoad()
		}

		macro.SymfonyCacheClear()
		macro.SymfonyAssetsInstall()
		macro.SymfonyImportMapInstall()

		macro.SoureCodeScreenStart()
		macro.SymfonyServerStart()

		shell.GetLogger().LogDuration()

		return nil
	},
}
