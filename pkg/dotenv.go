// Package dotenv provide methods to use on .env files.
package dotenv

import (
	"errors"
	"os"
	"reflect"
	"strings"
)

var (
	errReadingFile   = errors.New("can not read file")
	errFileIsEmpty   = errors.New(".env is empty or does not have key value pairs")
	errWrongFormat   = errors.New(".env is not in correct format")
	errAlreadyExists = errors.New("key value pair already exists")
	errMissingValue  = errors.New("value for the given key is not found")
)

type EnvContent struct {
	keyValuePairs map[string]string
}

// LoadFromString loads the content of .env file from multi-lined string.
func (env *EnvContent) LoadFromString(envContents string) (map[string]string, error) {

	lines := strings.Split(envContents, "\n")

	for _, line := range lines {

		if len(line) == 0 {
			continue
		} else if string(line[0]) == "#" {
			continue
		} else {
			s := strings.Split(line, "=")

			if s[0] == line {
				s = strings.Split(line, ":")
			}
			if s[0] == line || len(s) > 2 {
				return env.keyValuePairs, errWrongFormat
			}

			key, value := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
			env.keyValuePairs[key] = value
		}
	}

	return env.keyValuePairs, nil
}

// LoadFromFile loads the content of a given .env file
func (env *EnvContent) LoadFromFile(fileName string) (map[string]string, error) {

	env.keyValuePairs = make(map[string]string)
	emptyMap := make(map[string]string)

	err := error(nil)

	fileContent, err := os.ReadFile(fileName)

	if err != nil {
		return emptyMap, errReadingFile
	}

	_, err = env.LoadFromString(string(fileContent))

	if err != nil {
		return emptyMap, err
	}

	if reflect.DeepEqual(env.keyValuePairs, emptyMap) {
		return emptyMap, errFileIsEmpty
	}

	return env.keyValuePairs, err
}

// LoadFromFiles loads the content of given .env files
func (env *EnvContent) LoadFromFiles(fileNames []string) (map[string]string, error) {

	env.keyValuePairs = make(map[string]string)
	emptyMap := make(map[string]string)

	err := error(nil)

	for _, fileName := range fileNames {

		fileContent, err := os.ReadFile(fileName)

		if err != nil {
			err = errReadingFile
			continue
		}

		_, err = env.LoadFromString(string(fileContent))
	}

	if reflect.DeepEqual(env.keyValuePairs, emptyMap) {
		return emptyMap, errFileIsEmpty
	}
	return env.keyValuePairs, err
}

// GetEnv retrives the key value pairs of the .env files
func (env *EnvContent) GetEnv() map[string]string {
	return env.keyValuePairs
}

// Get retrives a value for a specific key
func (env *EnvContent) Get(key string) (string, error) {
	value := env.keyValuePairs[key]
	if value == "" {
		return value, errMissingValue
	}
	return value, nil
}
