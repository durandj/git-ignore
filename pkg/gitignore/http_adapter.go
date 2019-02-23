package gitignore

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
)

const defaultBaseURL string = "https://gitignore.io/api"

var sanitizers []*regexp.Regexp

func init() {
	sanitizers = []*regexp.Regexp{
		regexp.MustCompile("# Created by https://www.gitignore.io/api/[^\\s]+"),
		regexp.MustCompile("# Edit at https://www.gitignore.io/\\?templates=[^\\s]+"),
		regexp.MustCompile("# End of https://www.gitignore.io/api/[^\\s]+"),
	}
}

// HTTPAdapter is an adapter that can retrieve content by making HTTP
// requests to gitignore.io.
type HTTPAdapter struct {
	baseURL string
}

// List returns the list of options that can be used to generate a
// gitignore file.
func (adapter *HTTPAdapter) List() ([]string, error) {
	listURL := fmt.Sprintf("%s/list", adapter.baseURL)
	response, err := http.Get(listURL)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to get options list")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"HTTP error %d: %s",
			response.StatusCode,
			response.Status,
		)
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

	url := fmt.Sprintf(
		"%s/%s",
		adapter.baseURL,
		url.PathEscape(strings.Join(options, ",")),
	)
	response, err := http.Get(url)
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate gitignore file")
	}

	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"HTTP error %d: %s",
			response.StatusCode,
			response.Status,
		)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.Wrap(err, "Unable to generate gitignore file")
	}

	contents := string(body)
	for _, sanitizer := range sanitizers {
		contents = sanitizer.ReplaceAllLiteralString(contents, "")
	}

	return strings.TrimSpace(contents), nil
}

// Source generates potentially new gitignore data that can be
// stored or cached by other adapters.
func (adapter *HTTPAdapter) Source() (map[string]string, error) {
	options, err := adapter.List()
	if err != nil {
		return nil, errors.Wrap(err, "Unable to retrieve list of options")
	}

	mapping := make(map[string]string)
	for _, option := range options {
		content, err := adapter.Generate([]string{option})
		if err != nil {
			return nil, errors.Wrap(
				err,
				fmt.Sprintf("Unable to generate data for %s", option),
			)
		}

		mapping[option] = content
	}

	return mapping, nil
}

// NewHTTPAdapter creates a new HTTP adapter with the given base URL
// as the target host. If no URL is given then the default is used.
func NewHTTPAdapter(baseURL string) *HTTPAdapter {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &HTTPAdapter{
		baseURL: baseURL,
	}
}
