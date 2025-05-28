package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/SoureCode/kyx/commands"
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/SoureCode/kyx/tools"
	"github.com/pkg/errors"
	"github.com/symfony-cli/console"
)

var (
	// version is overridden at linking time
	version = "0.0.1"
	// channel is overridden at linking time
	channel = "dev"
	// overridden at linking time
	buildDate    = time.Now().Format(time.RFC3339)
	toolsMapping = tools.Mapping{
		"phpstan":                  "https://github.com/phpstan/phpstan/releases/latest/download/phpstan.phar",
		"php-cs-fixer":             "https://github.com/PHP-CS-Fixer/PHP-CS-Fixer/releases/latest/download/php-cs-fixer.phar",
		"infection":                "https://github.com/infection/infection/releases/latest/download/infection.phar",
		"composer-require-checker": "https://github.com/maglnet/ComposerRequireChecker/releases/latest/download/composer-require-checker.phar",
		"composer-normalize":       "https://github.com/ergebnis/composer-normalize/releases/latest/download/composer-normalize.phar",
		"composer-unused":          "https://github.com/composer-unused/composer-unused/releases/latest/download/composer-unused.phar",
	}
)

func main() {
	cmds := append(commands.GetCommands(), tools.GetCommands(toolsMapping)...)
	cmds = append(cmds, tools.CreateCommand("test"))
	toolNames := tools.GetNames(toolsMapping)

	args := os.Args

	if len(args) >= 2 && slices.Contains(toolNames, args[1]) {
		cmd := tools.NewToolCommand(args[1], toolsMapping, args[2:]...)

		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing command for tool %s: %v\n", args[1], err)
			os.Exit(cmd.ExitCode())
		}

		os.Exit(0)
	}

	if len(args) >= 2 && args[1] == "test" {
		p := project.GetProject()
		pd := p.GetDirectory()

		if p.HasDevDependency("symfony/phpunit-bridge") {
			cmd := shell.NewPHPCommand(append([]string{filepath.Join(pd, "vendor", "bin", "simple-phpunit")}, args[2:]...)...).WithPassthrough()

			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error executing simple-phpunit: %v\n", err)
				os.Exit(cmd.ExitCode())
			}

			os.Exit(0)
		} else if p.HasDevDependency("phpunit/phpunit") {
			cmd := shell.NewPHPCommand(append([]string{filepath.Join(pd, "vendor", "bin", "phpunit")}, args[2:]...)...).WithPassthrough()

			if err := cmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Error executing phpunit: %v\n", err)
				os.Exit(cmd.ExitCode())
			}

			os.Exit(0)
		}

		os.Exit(0)
	}

	app := &console.Application{
		Name:      "kyx",
		Usage:     "kyx [command] [options]",
		Copyright: fmt.Sprintf("(c) 2025-%d Jason Schilling", time.Now().Year()),
		Commands:  cmds,
		Action: func(ctx *console.Context) error {
			if ctx.Args().Len() == 0 {
				return commands.WelcomeAction(ctx)
			}

			return console.ShowAppHelpAction(ctx)
		},
		Version:   version,
		Channel:   channel,
		BuildDate: buildDate,
	}

	err := app.Run(args)

	if err != nil {
		panic(errors.Wrap(err, "could not run application"))
	}
}
