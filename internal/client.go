package internal

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Client is an object used to interact with the gitignore provider.
// It knows how to retrieve a list of supported apps to be ignored
// and turn them into a gitignore file.
type Client struct {
	Adapters []Adapter
}

// NewClient creates a new client for generating gitignore files.
func NewClient() (*Client, error) {
	gitAdapter, err := NewGitAdapter()
	if err != nil {
		return nil, fmt.Errorf("unable to create Git adapter: %w", err)
	}

	return &Client{
		Adapters: []Adapter{
			gitAdapter,
		},
	}, nil
}

// List returns a list of valid options for generating a gitignore
// file. Each of these options maps to a service or application that
// generates file that should be excluded from a git repository.
func (client *Client) List() ([]string, error) {
	adapterErrors := []error{}

	for _, adapter := range client.Adapters {
		options, err := adapter.List()

		if err != nil {
			adapterErrors = append(adapterErrors, err)

			continue
		}

		slices.SortFunc(options, func(a, b string) int {
			return strings.Compare(strings.ToLower(a), strings.ToLower(b))
		})

		return options, nil
	}

	return nil, fmt.Errorf("unable to retrieve option list:\n%s", adapterErrors)
}

// Generate generates a .gitignore file that excludes files based on
// the given options.
func (client *Client) Generate(options []string) (string, error) {
	if len(options) == 0 {
		return "", errors.New("must give at least one option")
	}

	adapterErrors := []error{}

	for _, adapter := range client.Adapters {
		validOptions, err := adapter.List()
		if err != nil {
			adapterErrors = append(adapterErrors, err)

			continue
		}

		for _, option := range options {
			if !slices.Contains(validOptions, option) {
				return "", fmt.Errorf("invalid option \"%s\"", option)
			}
		}

		content, err := adapter.Generate(options)
		if err != nil {
			adapterErrors = append(adapterErrors, err)

			continue
		}

		return content, nil
	}

	return "", fmt.Errorf("unable to generate gitignore:\n%s", adapterErrors)
}

// Update updates all local cache adapters.
func (client *Client) Update() error {
	for _, adapter := range client.Adapters {
		err := adapter.Update()
		if err != nil {
			return fmt.Errorf("unable to update adapter: %w", err)
		}
	}

	return nil
}
