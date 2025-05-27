package commands

import (
	"github.com/symfony-cli/console"
	"github.com/symfony-cli/terminal"
)

func WelcomeAction(c *console.Context) error {
	console.ShowVersion(c)
	terminal.Println(c.App.Usage)
	terminal.Printf("Show all commands with <info>%s help</>,\n", c.App.HelpName)
	terminal.Printf("Get help for a specific command with <info>%s help COMMAND</>.\n", c.App.HelpName)
	return nil
}

func GetCommands() []*console.Command {
	return []*console.Command{
		startCommand,
		stopCommand,
		resetCommand,
		deploymentCommand,
	}
}
