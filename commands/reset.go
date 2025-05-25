package commands

import (
	"github.com/SoureCode/kyx/macro"
	"github.com/SoureCode/kyx/shell"
	"github.com/symfony-cli/console"
)

var resetCommand = &console.Command{
	Name:        "reset",
	Description: "Resets application state",
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
		macro.SymfonyWorkerStop()
		macro.SoureCodeScreenStop()
		macro.WaitForDatabase()
		macro.SymfonyDoctrineDatabaseDrop()
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

		shell.GetLogger().LogDuration()

		return nil
	},
}
