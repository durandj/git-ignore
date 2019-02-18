package gitignore

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const baseURL string = "https://gitignore.io/api"

// Client is an object used to interact with the gitignore provider.
// It knows how to retrieve a list of supported apps to be ignored
// and turn them into a gitignore file.
type Client struct {
}

// List returns a list of valid options for generating a gitignore
// file. Each of these options maps to a service or application that
// generates file that should be excluded from a git repository.
func (client *Client) List() ([]string, error) {
	listURL := fmt.Sprintf("%s/list", baseURL)
	response, err := http.Get(listURL)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get options list")
	}

	defer response.Body.Close()

	options := make([]string, 0)

	// TODO(durandj): this could be done by using a custom splitfunc
	scanner := bufio.NewScanner(response.Body)
	for scanner.Scan() {
		line := scanner.Text()

		for _, option := range strings.Split(line, ",") {
			option = strings.TrimSpace(option)
			if option == "" {
				continue
			}

			options = append(options, option)
		}
	}

	return options, nil
}

// Generate generates a .gitignore file that excludes files based on
// the given options.
func (client *Client) Generate(options []string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("Must give at least one option")
	}

	validOptions, err := client.List()
	if err != nil {
		return "", errors.Wrap(err, "Unable to validate options")
	}

	for _, option := range options {
		if !includes(validOptions, option) {
			return "", fmt.Errorf("Invalid option \"%s\"", option)
		}
	}

	url := fmt.Sprintf("%s/%s", baseURL, url.PathEscape(strings.Join(options, ",")))
	response, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate gitignore file")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate gitignore file")
	}

	return string(body), nil
}

func includes(expected []string, value string) bool {
	for _, valid := range expected {
		if valid == value {
			return true
		}
	}

	return false
}
