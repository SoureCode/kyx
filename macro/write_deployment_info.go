package macro

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/SoureCode/kyx/internal/php"
	"github.com/SoureCode/kyx/project"
	"github.com/SoureCode/kyx/shell"
)

func WriteDeploymentInfo(repo string) {
	p := project.GetProject()
	dir := p.GetDirectory()
	logger := shell.GetLogger()

	if !filepath.IsAbs(repo) {
		repo = filepath.Clean(filepath.Join(dir, repo))
	}

	deploymentFile := filepath.Join(dir, "deployment.php")
	deploymentInfo, err := php.LoadFileDump(deploymentFile)

	if err != nil {
		logger.Errorln("Error loading deployment info:", err)
		// fail as the deployment info is required
		os.Exit(1)
	}

	cmd := shell.NewGitCommand("-C", repo, "rev-parse", "HEAD").
		WithLogger(logger).
		WithLogLevel(0)

	if err := cmd.Run(); err != nil {
		logger.Errorln("Error getting git commit hash:", err)
		os.Exit(1)
	}

	deploymentInfo["git_commit"] = strings.TrimSpace(cmd.Stdout())

	if err := php.WriteFileDump(deploymentFile, deploymentInfo); err != nil {
		logger.Errorln("Error writing deployment info:", err)
		os.Exit(1)
	} else {
		logger.Infoln("Deployment info written to", deploymentFile)
	}
}
