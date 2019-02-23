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

// Source is any type that can produce new gitignore data instead of
// just caching that data.
type Source interface {
	// Source generates potentially new gitignore data that can be
	// stored or cached by other adapters.
	Source() (map[string]string, error)
}

// Cache stores gitignore data but cannot produce new or novel data.
type Cache interface {
	Cache(ignoreMapping map[string]string) error
}
