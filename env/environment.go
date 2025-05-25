package env

import (
	"os"
)

type Environment struct {
	variables map[string]string
}

func (e *Environment) Lookup(key string) (string, bool) {
	value, exists := e.variables[key]

	if !exists {
		value, exists = os.LookupEnv(key)

		if !exists {
			return "", false
		}
	}

	return value, exists
}

func (e *Environment) Get(key string) string {
	value, _ := e.Lookup(key)
	return value
}

func (e *Environment) IsProd() bool {
	prod, exists := e.Lookup("APP_ENV")

	if !exists {
		return false
	}

	return prod == "prod"
}

func (e *Environment) IsDev() bool {
	dev, exists := e.Lookup("APP_ENV")

	if !exists {
		return false
	}

	return dev == "dev"
}

func (e *Environment) Environ() []string {
	env := []string{}

	for key, value := range e.variables {
		env = append(env, key+"="+value)
	}

	return env
}

func NewEnvironment(directory string) (*Environment, error) {
	envMap := map[string]string{}

	err := loadFile(envMap, directory, ".env")

	if err != nil {
		return nil, err
	}

	err = loadFile(envMap, directory, ".env.local")

	if err != nil {
		return nil, err
	}

	envName := os.Getenv("APP_ENV")

	if envName == "" {
		if value, exists := envMap["APP_ENV"]; exists {
			envName = value
		}
	}

	if envName != "" {
		err = loadFile(envMap, directory, ".env."+envName)

		if err != nil {
			return nil, err
		}

		err = loadFile(envMap, directory, ".env."+envName+".local")

		if err != nil {
			return nil, err
		}
	}

	err = loadPHPFile(envMap, directory, ".env.local.php")

	if err != nil {
		return nil, err
	}

	return &Environment{
		variables: envMap,
	}, nil
}
