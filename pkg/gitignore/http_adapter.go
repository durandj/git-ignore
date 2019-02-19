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

// HTTPAdapter is an adapter that can retrieve content by making HTTP
// requests to gitignore.io.
type HTTPAdapter struct {
}

// List returns the list of options that can be used to generate a
// gitignore file.
func (adapter *HTTPAdapter) List() ([]string, error) {
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

// Generate creates a gitignore file with the given options.
func (adapter *HTTPAdapter) Generate(options []string) (string, error) {
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
