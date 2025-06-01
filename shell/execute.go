package shell

import (
	"os/exec"
	"path/filepath"

	"github.com/SoureCode/kyx/project"
	"github.com/pkg/errors"
)

var (
	SymfonyNotFoundError = errors.New("symfony binary not found")
)

func NewSymfonyCommand(args ...string) (*Command, error) {
	binary, err := exec.LookPath("symfony")

	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return nil, SymfonyNotFoundError
		}

		return nil, errors.Wrap(err, "could not find symfony binary")
	}

	cmd := NewCommand(binary).WithArgs(args...)

	return cmd, nil
}

func NewPHPCommand(args ...string) *Command {
	if symfonyCommand, err := NewSymfonyCommand(append([]string{"php"}, args...)...); err == nil {
		return symfonyCommand
	} else if !errors.Is(err, SymfonyNotFoundError) {
		panic(errors.Wrap(err, "could not create symfony command"))
	}

	binary, err := exec.LookPath("php")

	if err != nil {
		panic(errors.Wrap(err, "could not find php binary"))
	}

	return NewCommand(binary).WithArgs(args...)
}

func NewConsoleCommand(args ...string) *Command {
	if symfonyCommand, err := NewSymfonyCommand(append([]string{"console"}, args...)...); err == nil {
		return symfonyCommand
	} else if !errors.Is(err, SymfonyNotFoundError) {
		panic(errors.Wrap(err, "could not create symfony command"))
	}

	p := project.GetProject()
	pd := p.GetDirectory()

	consoleFilePath := filepath.Join(pd, "bin", "console")

	return NewPHPCommand(append([]string{consoleFilePath}, args...)...)
}

func NewComposerCommand(args ...string) *Command {
	if symfonyCommand, err := NewSymfonyCommand(append([]string{"composer"}, args...)...); err == nil {
		return symfonyCommand
	} else if !errors.Is(err, SymfonyNotFoundError) {
		panic(errors.Wrap(err, "could not create symfony command"))
	}

	binary, err := exec.LookPath("composer")

	if err != nil {
		panic(errors.Wrap(err, "could not find composer binary"))
	}

	return NewCommand(binary).WithArgs(args...)
}

func NewDockerCommand(args ...string) *Command {
	binary, err := exec.LookPath("docker")

	if err != nil {
		panic(errors.Wrap(err, "could not find docker binary"))
	}

	return NewCommand(binary).WithArgs(args...)
}
