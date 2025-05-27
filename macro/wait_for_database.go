package macro

import (
	"time"

	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"github.com/pkg/errors"
)

func WaitForDatabase() {
	logger := shell.GetLogger()
	p := project.GetProject()

	if p.HasDependency("doctrine/doctrine-bundle") {
		retryTimeout := time.Millisecond * 100
		attemptsLeftToReachDatabase := 100 // 100 attempts with 100ms interval = 10 seconds
		var databaseError error = nil
		exitCode := 0

		logger.Log("Waiting for database to be ready")

		for attemptsLeftToReachDatabase > 0 {
			cmd := shell.NewConsoleCommand("dbal:run-sql", "-q", "SELECT 1")

			if err := cmd.Run(); err == nil {
				logger.Logln()
				logger.Logln("The database is now ready and reachable")
				return
			} else {
				databaseError = err
				exitCode = cmd.ExitCode()

				if exitCode == 255 {
					logger.Logln()
					panic(errors.Wrap(databaseError, "unrecoverable error encountered while checking database connection"))
				}
			}

			logger.Log(".")
			attemptsLeftToReachDatabase--
			time.Sleep(retryTimeout)
		}

		logger.Logln()
		panic(errors.Wrap(databaseError, "he database is not up or not reachable"))
	}
}
