package gitignore

import (
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// FSAdapter is an adapter for git-ignore that reads from the file
// system.
type FSAdapter struct {
	baseDir string
}

// List returns the list of options that can be used to generate a
// gitignore file.
func (adapter *FSAdapter) List() ([]string, error) {
	files, err := ioutil.ReadDir(adapter.baseDir)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get list from filesystem")
	}

	options := make([]string, len(files))
	for index, file := range files {
		options[index] = file.Name()
	}

	return options, nil
}

// Generate creates a gitignore file with the given options.
func (adapter *FSAdapter) Generate(options []string) (string, error) {
	var builder strings.Builder

	for _, option := range options {
		contents, err := ioutil.ReadFile(path.Join(adapter.baseDir, option))

		if err != nil {
			return "", errors.Wrap(err, "Unable to read gitignore recipe")
		}

		builder.Write(contents)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

// Cache stores gitignore data but cannot produce new or novel data.
func (adapter *FSAdapter) Cache(ignoreMapping map[string]string) error {
	err := os.MkdirAll(adapter.baseDir, os.ModePerm)
	if err != nil {
		return errors.Wrap(err, "Unable to create directory for cache")
	}

	for option, contents := range ignoreMapping {
		filePath := path.Join(adapter.baseDir, option)

		err := ioutil.WriteFile(filePath, []byte(contents), 0644)
		if err != nil {
			return errors.Wrap(err, "Unable to write cache file")
		}
	}

	return nil
}

// NewFSAdapter creates an adapter that is able to read from the file
// system to find data to use for generating gitignore files.
func NewFSAdapter(baseDir string) (*FSAdapter, error) {
	if baseDir == "" {
		currentUser, err := user.Current()
		if err != nil {
			return nil, errors.Wrap(err, "Unable to get user home directory")
		}

		baseDir = path.Join(currentUser.HomeDir, ".local", "share", "git-ignore")
	}

	return &FSAdapter{baseDir: baseDir}, nil
}
