package project

import (
	"encoding/json"
	"os"
)

type composerJson struct {
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

func (cj *composerJson) GetDependencies() []string {
	var dependencies []string

	for key := range cj.Require {
		dependencies = append(dependencies, key)
	}

	return dependencies
}

func (cj *composerJson) GetDevDependencies() []string {
	var devDependencies []string

	for key := range cj.RequireDev {
		devDependencies = append(devDependencies, key)
	}

	return devDependencies
}

func loadComposerJson(path string) (*composerJson, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, os.ErrNotExist
	}

	file, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var composerData composerJson
	if err := json.Unmarshal(file, &composerData); err != nil {
		return nil, err
	}

	return &composerData, nil
}
