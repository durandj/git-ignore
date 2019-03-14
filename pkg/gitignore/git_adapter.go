package gitignore

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/src-d/go-git.v4"
)

// DefaultGitRepo is the default repository to use for gitignore files.
const DefaultGitRepo string = "https://github.com/github/gitignore.git"

// GitAdapter is an adapter for pulling gitignore data from a git
// repository.
type GitAdapter struct {
	RepoDirectory string
	RepoURL       string
}

// List returns the list of options that can be used to generate a
// gitignore file.
func (adapter *GitAdapter) List() ([]string, error) {
	options := []string{}

	err := filepath.Walk(adapter.RepoDirectory, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Unable to file gitignore files")
		}

		if path.Ext(filePath) != ".gitignore" {
			return nil
		}

		option := path.Base(strings.Replace(filePath, ".gitignore", "", 1))
		options = append(options, option)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return options, nil
}

// Generate creates a gitignore file with the given options.
func (adapter *GitAdapter) Generate(options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("Must give at least one option")
	}

	validOptions, err := adapter.List()
	if err != nil {
		return "", errors.Wrap(err, "Unable to validate options")
	}

	for _, option := range options {
		if !includes(validOptions, option) {
			return "", fmt.Errorf("Invalid option \"%s\"", option)
		}
	}

	var builder strings.Builder
	for _, option := range options {
		filePath, err := adapter.findOptionFile(option)
		if err != nil {
			return "", errors.Wrap(err, "Unable to find file")
		}

		contents, err := ioutil.ReadFile(filePath)
		if err != nil {
			message := fmt.Sprintf("Unable to read gitignore data for %s", option)
			return "", errors.Wrap(err, message)
		}

		builder.WriteString(fmt.Sprintf("### %s ###\n", option))
		builder.Write(contents)
		builder.WriteString("\n")
	}

	return builder.String(), nil
}

// Update updates this plugin's local data.
func (adapter *GitAdapter) Update() error {
	_, err := os.Stat(adapter.RepoDirectory)
	repoExists := err != nil

	if repoExists {
		_, err = git.PlainClone(adapter.RepoDirectory, false, &git.CloneOptions{
			URL: adapter.RepoURL,
		})

		return errors.Wrap(err, "Unable to clone repository")
	}

	repository, err := git.PlainOpen(adapter.RepoDirectory)
	if err != nil {
		return errors.Wrap(err, "Unable to open repository")
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return errors.Wrap(err, "Unable to get working tree")
	}

	err = worktree.Pull(&git.PullOptions{})
	if err == git.NoErrAlreadyUpToDate {
		err = nil
	}

	return err
}

func (adapter *GitAdapter) findOptionFile(option string) (string, error) {
	filename := option + ".gitignore"
	filePath := ""

	err := filepath.Walk(adapter.RepoDirectory, func(currentFile string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Error while looking for file")
		}

		if strings.HasSuffix(currentFile, filename) {
			filePath = currentFile

			return io.EOF
		}

		return nil
	})

	if err != nil && err != io.EOF {
		message := fmt.Sprintf("Unable to find file for %s", option)

		return "", errors.Wrap(err, message)
	}

	if filePath == "" {
		return "", errors.Errorf("Unable to find file for %s", option)
	}

	return filePath, nil
}

// NewGitAdapter creates a new adapter for working with Git
// repositories.
func NewGitAdapter() (*GitAdapter, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get user home directory")
	}

	userHome := currentUser.HomeDir

	return &GitAdapter{
		RepoDirectory: path.Join(userHome, ".local", "share", "git-ignore", "gitignore"),
		RepoURL:       DefaultGitRepo,
	}, nil
}
