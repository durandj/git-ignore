package internal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5"
)

// DefaultGitRepo is the default repository to use for gitignore files.
const DefaultGitRepo string = "https://github.com/github/gitignore.git"

// GitAdapter is an adapter for pulling gitignore data from a git
// repository.
type GitAdapter struct {
	RepoDirectory string
	RepoURL       string
}

// NewGitAdapter creates a new adapter for working with Git
// repositories.
func NewGitAdapter() (*GitAdapter, error) {
	currentUser, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("unable to get user home directory: %w", err)
	}

	userHome := currentUser.HomeDir

	return &GitAdapter{
		RepoDirectory: path.Join(userHome, ".local", "share", "git-ignore", "gitignore"),
		RepoURL:       DefaultGitRepo,
	}, nil
}

// List returns the list of options that can be used to generate a
// gitignore file.
func (adapter *GitAdapter) List() ([]string, error) {
	options := []string{}

	err := filepath.Walk(adapter.RepoDirectory, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("unable to file gitignore files: %w", err)
		}

		if path.Ext(filePath) != ".gitignore" {
			return nil
		}

		option := path.Base(strings.Replace(filePath, ".gitignore", "", 1))
		options = append(options, option)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to read gitignore repository: %w", err)
	}

	return options, nil
}

// Generate creates a gitignore file with the given options.
func (adapter *GitAdapter) Generate(options []string) (string, error) {
	if len(options) == 0 {
		return "", errors.New("must give at least one option")
	}

	validOptions, err := adapter.List()
	if err != nil {
		return "", fmt.Errorf("unable to validate options for generating ignore file: %w", err)
	}

	for _, option := range options {
		if !slices.Contains(validOptions, option) {
			return "", fmt.Errorf("invalid option \"%s\"", option)
		}
	}

	var builder strings.Builder
	for _, option := range options {
		filePath, err := adapter.findOptionFile(option)
		if err != nil {
			return "", fmt.Errorf("unable to find file: %w", err)
		}

		contents, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("unable to read gitignore data for %s: %w", option, err)
		}

		if !bytes.HasPrefix(contents, []byte("###")) {
			builder.WriteString(fmt.Sprintf("### %s ###\n", option))
		}

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
		//nolint:exhaustruct // Only URL is required
		_, err = git.PlainClone(adapter.RepoDirectory, false, &git.CloneOptions{
			URL: adapter.RepoURL,
		})

		if err != nil {
			return fmt.Errorf("unable to clone repository: %w", err)
		}

		return nil
	}

	repository, err := git.PlainOpen(adapter.RepoDirectory)
	if err != nil {
		return fmt.Errorf("unable to open repository: %w", err)
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return fmt.Errorf("unable to get working tree: %w", err)
	}

	//nolint:exhaustruct // all fields are optional
	err = worktree.Pull(&git.PullOptions{})
	if errors.Is(err, git.NoErrAlreadyUpToDate) {
		err = nil
	}

	if err != nil {
		return fmt.Errorf("unable to update gitignore repository: %w", err)
	}

	return nil
}

func (adapter *GitAdapter) findOptionFile(option string) (string, error) {
	filename := option + ".gitignore"
	filePath := ""

	err := filepath.Walk(adapter.RepoDirectory, func(currentFile string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error while looking for file: %w", err)
		}

		if strings.HasSuffix(currentFile, filename) {
			filePath = currentFile

			return io.EOF
		}

		return nil
	})

	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("unable to find file for %s: %w", option, err)
	}

	if filePath == "" {
		return "", fmt.Errorf("unable to find file for %s", option)
	}

	return filePath, nil
}
