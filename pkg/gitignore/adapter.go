package gitignore

// Adapter is any adapter that the git-ignore client can use to
// retrieve content for generating a gitignore file.
type Adapter interface {
	// List returns the list of options that can be used to generate a
	// gitignore file.
	List() ([]string, error)

	// Generate creates a gitignore file with the given options.
	Generate(options []string) (string, error)
}
