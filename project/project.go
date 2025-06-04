package project

import (
	"path/filepath"
	"slices"

	"github.com/SoureCode/kyx/env"
	"github.com/pkg/errors"
)

var singletonProject *Project

func GetProject() *Project {
	if singletonProject == nil {
		directory, err := getDirectory()

		if err != nil {
			panic(errors.Wrap(err, "failed to get project directory"))
		}

		project, err := newProject(directory)

		if err != nil {
			panic(errors.Wrap(err, "failed to create project"))
		}

		singletonProject = project
	}

	return singletonProject
}

type Project struct {
	directory       string
	dependencies    []string
	devDependencies []string
	environment     *env.Environment
}

func newProject(directory string) (*Project, error) {
	composerPath := filepath.Join(directory, "composer.json")
	composerData, err := loadComposerJson(composerPath)

	if err != nil {
		return nil, err
	}

	environment, err := env.NewEnvironment(directory)

	if err != nil {
		return nil, errors.Wrap(err, "failed to load environment")
	}

	return &Project{
		directory:       directory,
		dependencies:    composerData.GetDependencies(),
		devDependencies: composerData.GetDevDependencies(),
		environment:     environment,
	}, nil
}

func HasProject() bool {
	directory, err := getDirectory()

	if err != nil {
		return false
	}

	composerPath := filepath.Join(directory, "composer.json")
	_, err = loadComposerJson(composerPath)

	return err == nil
}

func (p *Project) HasDependency(dependency string) bool {
	return slices.Contains(p.dependencies, dependency)
}

func (p *Project) HasDevDependency(dependency string) bool {
	return slices.Contains(p.devDependencies, dependency)
}

func (p *Project) GetDirectory() string {
	return p.directory
}

func (p *Project) GetEnvironment() *env.Environment {
	return p.environment
}
