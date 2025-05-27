package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/SoureCode/kyx/project"
	"github.com/pkg/errors"
)

type Command struct {
	binary      string
	args        []string
	stdin       io.Reader
	stdout      *OutputWriter
	stderr      *OutputWriter
	project     *project.Project
	exitCode    int
	executed    bool
	logger      *Logger
	print       bool
	passthrough bool
}

func (e *Command) Execute() {
	if e.executed {
		panic(errors.New("execution already executed"))
	}

	e.executed = true

	if e.logger != nil {
		e.logger.Logln("Executing command: ", e.binary, " ", strings.Join(e.args, " "))
	}

	go func() {
		for line := range e.stdout.Lines {
			if e.print {
				if e.logger != nil {
					e.logger.Logln(" >" + line)
				} else {
					fmt.Println(line)
				}
			}
		}
	}()

	go func() {
		for line := range e.stderr.Lines {
			if e.print {
				if e.logger != nil {
					e.logger.Logln(" >" + line)
				} else {
					fmt.Println(line)
				}
			}
		}
	}()

	cmd := exec.Command(e.binary, e.args...)
	cmd.Stdin = e.stdin
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	cmd.Dir = e.project.GetDirectory()

	if e.passthrough {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}

	env := e.project.GetEnvironment()
	environ := append(os.Environ(), env.Environ()...)
	cmd.Env = append(cmd.Env, environ...)

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		e.exitCode = 1
		return
	}

	waitChannel := make(chan error, 1)
	go func() {
		waitChannel <- cmd.Wait()
		close(waitChannel)
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel)
	defer signal.Stop(signalChannel)

	for {
		select {
		case err := <-waitChannel:
			e.exitCode = 0

			if err == nil {
				e.stdout.CloseLines()
				e.stderr.CloseLines()
				return
			}

			if !strings.Contains(err.Error(), "exit status") {
				fmt.Fprintln(os.Stderr, err)
			}

			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					e.exitCode = status.ExitStatus()
				}
			}

			return
		case sig := <-signalChannel:
			// this one in particular should be skipped as we don't want to
			// send it back to child because it's about it
			if sig == syscall.SIGCHLD {
				continue
			}

			if cmd.Process != nil {
				if err := cmd.Process.Signal(sig); err != nil && !strings.Contains(err.Error(), "process already finished") {
					fmt.Fprintln(os.Stderr, "error sending signal", sig, err)
				}
			}
		}
	}
}

func NewCommand(binary string) *Command {
	return &Command{
		binary:   binary,
		stdin:    os.Stdin,
		stdout:   NewOutputWriter(),
		stderr:   NewOutputWriter(),
		project:  project.GetProject(),
		exitCode: -1,
		executed: false,
		logger:   nil,
		print:    false,
	}
}

func (e *Command) WithPassthrough() *Command {
	e.passthrough = true
	return e
}

func (e *Command) WithArgs(args ...string) *Command {
	e.args = args
	return e
}

func (e *Command) WithLogger(logger *Logger) *Command {
	e.logger = logger
	return e
}

func (e *Command) WithProject(p *project.Project) *Command {
	e.project = p
	return e
}

func (e *Command) WithPrintOutput(p bool) *Command {
	e.print = p
	return e
}

func (e *Command) WithStdin(r io.Reader) *Command {
	e.stdin = r
	return e
}

func (e *Command) ExitCode() int {
	return e.exitCode
}

func (e *Command) Stdout() string {
	return e.stdout.String()
}

func (e *Command) Stderr() string {
	return e.stderr.String()
}

func (e *Command) Run() error {
	e.Execute()

	if e.ExitCode() != 0 {
		return errors.Wrapf(
			errors.New(e.Stderr()),
			"command '%s %s' failed with exit code %d",
			e.binary,
			strings.Join(e.args, " "),
			e.ExitCode(),
		)
	}

	return nil
}
