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
	logger      *BaseLogger
	logLevel    int // 1: error, 2: warning, 3: info, 4: debug (default), 5: trace
	passthrough bool
}

func (e *Command) Execute() error {
	if e.executed {
		panic(errors.New("execution already executed"))
	}

	e.executed = true

	if e.logger != nil {
		e.logger.Logln("Executing command: ", e.binary, " ", strings.Join(e.args, " "))
	}

	cmd := exec.Command(e.binary, e.args...)
	cmd.Stdin = e.stdin
	cmd.Stdout = e.stdout
	cmd.Stderr = e.stderr
	cmd.Dir = e.project.GetDirectory()

	if e.passthrough {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	} else {
		go func() {
			for line := range e.stdout.Lines {
				if e.logger != nil {
					e.logger.doLogln(e.logLevel, " >"+line)
				}
			}
		}()

		go func() {
			for line := range e.stderr.Lines {
				if e.logger != nil {
					e.logger.doLogln(e.logLevel, " >"+line)
				}
			}
		}()
	}

	env := e.project.GetEnvironment()
	environ := append(os.Environ(), env.Environ()...)
	cmd.Env = append(cmd.Env, environ...)

	if err := cmd.Start(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		e.exitCode = 1
		return errors.Wrap(err, "failed to start command")
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
				return nil
			}

			if !strings.Contains(err.Error(), "exit status") {
				fmt.Fprintln(os.Stderr, err)
			}

			if exitErr, ok := err.(*exec.ExitError); ok {
				if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
					e.exitCode = status.ExitStatus()
				}
			}

			return err
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
		logLevel: 3, // default log level is info
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

func (e *Command) WithLogger(logger *BaseLogger) *Command {
	e.logger = logger

	return e
}

func (e *Command) WithProject(p *project.Project) *Command {
	e.project = p
	return e
}

// WithLogLevel sets the log level for the command execution.
// 1: error
// 2: warning
// 3: info (default)
// 4: debug
// 5: trace
func (e *Command) WithLogLevel(logLevel int) *Command {
	e.logLevel = logLevel
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
	err := e.Execute()

	if err != nil {
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
