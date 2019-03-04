package gitignore

import (
	"fmt"

	"github.com/pkg/errors"
)

// Client is an object used to interact with the gitignore provider.
// It knows how to retrieve a list of supported apps to be ignored
// and turn them into a gitignore file.
type Client struct {
	Adapters []Adapter
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

		return options, nil
	}

	return nil, fmt.Errorf("Unable to retrieve option list:\n%s", adapterErrors)
}

// Generate generates a .gitignore file that excludes files based on
// the given options.
func (client *Client) Generate(options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("Must give at least one option")
	}

	adapterErrors := []error{}

	for _, adapter := range client.Adapters {
		validOptions, err := adapter.List()
		if err != nil {
			adapterErrors = append(adapterErrors, err)
			continue
		}

		for _, option := range options {
			if !includes(validOptions, option) {
				return "", fmt.Errorf("Invalid option \"%s\"", option)
			}
		}

		content, err := adapter.Generate(options)
		if err != nil {
			adapterErrors = append(adapterErrors, err)
			continue
		}

		return content, nil
	}

	return "", fmt.Errorf("Unable to generate gitignore:\n%s", adapterErrors)
}

// Update updates all local cache adapters.
func (client *Client) Update() error {
	mappings := make(map[string]string)
	for _, adapter := range client.Adapters {
		sourceAdapter, ok := adapter.(Source)
		if !ok {
			continue
		}

		adapterMappings, err := sourceAdapter.Source()
		if err != nil {
			return errors.Wrap(err, "Unable to collect ignore mappings")
		}

		for key, value := range adapterMappings {
			mappings[key] = value
		}
	}

	atLeastOneSuccess := false
	cacheErrors := []error{}
	for _, adapter := range client.Adapters {
		cacheAdapter, ok := adapter.(Cache)
		if !ok {
			continue
		}

		err := cacheAdapter.Cache(mappings)
		if err != nil {
			cacheErrors = append(cacheErrors, err)
			continue
		}

		atLeastOneSuccess = true
	}

	if len(cacheErrors) > 0 && !atLeastOneSuccess {
		return fmt.Errorf("Unable to update, all adapters failed\n%s", cacheErrors)
	}

	return nil
}

// NewClient creates a new client for generating gitignore files.
func NewClient() (*Client, error) {
	fsAdapter, err := NewFSAdapter("")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to create filesystem adapter")
	}

	return &Client{
		Adapters: []Adapter{
			fsAdapter,
			NewHTTPAdapter(""),
		},
	}, nil
}

func includes(expected []string, value string) bool {
	for _, valid := range expected {
		if valid == value {
			return true
		}
	}

	return false
}
