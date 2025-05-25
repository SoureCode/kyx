package tools

import (
	"github.com/pkg/errors"
	"github.com/symfony-cli/console"
)

func CreateCommand(name string) *console.Command {
	return &console.Command{
		Hidden: console.Hide,
		// we use an alias to avoid the command being shown in the help but
		// still be available for completion
		Aliases:     []*console.Alias{{Name: name}},
		Description: "Run " + name,
		Action: func(c *console.Context) error {
			return console.IncorrectUsageError{
				ParentError: errors.New(`This command can only be run as "` + c.App.HelpName + ` ` + name + `"`),
			}
		},
	}
}

func GetCommands(mapping Mapping) []*console.Command {
	commands := make([]*console.Command, 0)

	for name, _ := range mapping {
		commands = append(commands, CreateCommand(name))
	}

	return commands
}
