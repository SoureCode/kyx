package macro

import (
	"fmt"
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func SentryDeploysNew(repo string, startedAt, stoppedAt time.Time) {
	p := project.GetProject()
	dir := p.GetDirectory()
	logger := shell.GetLogger()

	if !filepath.IsAbs(repo) {
		repo = filepath.Clean(filepath.Join(dir, repo))
	}

	cmd := shell.NewGitCommand("-C", repo, "rev-parse", "HEAD").
		WithLogger(logger).
		WithLogLevel(3)

	if err := cmd.Run(); err != nil {
		logger.Errorln("Error getting git commit hash:", err)
		os.Exit(1)
	}

	hash := strings.TrimSpace(cmd.Stdout())
	env := p.GetEnvironment()

	// Only run if Sentry is configured
	sentryUrl, ok := env.Lookup("SENTRY_URL")
	if !ok {
		logger.Infoln("Sentry url not configured, skipping sentry integration")
		return
	}

	if strings.HasSuffix(sentryUrl, "/") {
		sentryUrl = strings.TrimSuffix(sentryUrl, "/")
	}

	sentryOrg, ok := env.Lookup("SENTRY_ORG")
	if !ok {
		logger.Infoln("Sentry organization not configured, skipping sentry integration")
		return
	}

	sentryProject, ok := env.Lookup("SENTRY_PROJECT")
	if !ok {
		logger.Infoln("Sentry project not configured, skipping sentry integration")
		return
	}

	_, ok = env.Lookup("SENTRY_AUTH_TOKEN")
	if !ok {
		logger.Infoln("Sentry auth token not configured, skipping sentry integration")
		return
	}

	appEnv, ok := env.Lookup("APP_ENV")
	if !ok {
		logger.Infoln("APP_ENV not configured, skipping sentry integration")
		return
	}

	cmd = shell.NewSentryCommand(
		"--url="+sentryUrl,
		"releases", "list",
		"--raw",
		"--org="+sentryOrg,
		"--project="+sentryProject,
	).
		WithLogger(logger).
		WithLogLevel(0)

	if err := cmd.Run(); err != nil {
		logger.Errorln("Error listing Sentry releases:", err)
		os.Exit(1)
	}

	releases := strings.TrimSpace(cmd.Stdout())

	if !strings.Contains(releases, hash) {
		cmd = shell.NewSentryCommand(
			"releases", "new",
			"--url="+sentryUrl,
			"--org="+sentryOrg,
			"--project="+sentryProject,
			hash,
		).
			WithLogger(logger).
			WithLogLevel(0)

		if err := cmd.Run(); err != nil {
			logger.Errorln("Error creating new Sentry release:", err)
			os.Exit(1)
		}
	}

	cmd = shell.NewSentryCommand(
		"--url="+sentryUrl,
		"deploys", "new",
		// "--url="+filepath.Base(filepath.Dir(dir)), // @todo strategy to get the URL
		"--org="+sentryOrg,
		"--project="+sentryProject,
		"--started="+startedAt.Format(time.RFC3339),
		"--finished="+stoppedAt.Format(time.RFC3339),
		"--release="+hash,
		"--env="+appEnv,
	).
		WithLogger(logger).
		WithLogLevel(0)

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing sentry-cli: %v\n", err)
		os.Exit(cmd.ExitCode())
	}
}
