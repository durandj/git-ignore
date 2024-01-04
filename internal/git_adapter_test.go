package internal_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/durandj/git-ignore/internal"
)

func TestGitAdapterListShouldRetrieveAListOfOptions(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	options, err := adapter.List()
	require.NoError(t, err)

	require.Contains(t, options, "C")
	require.Contains(t, options, "C++")
}

func TestGitAdapterListShouldReturnAnErrorIfTheRepositoryDoesNotExist(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	_, err = adapter.List()

	require.Error(t, err)
}

func TestGitAdapterGenerateShouldReturnAnErrorWhenNoOptionsAreGiven(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	_, err = adapter.Generate([]string{})

	require.Error(t, err)
}

func TestGitAdapterGenerateShouldReturnAnErrorWhenTheRepositoryDoesNotExist(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	_ = os.RemoveAll(testDir)

	_, err = adapter.Generate([]string{"C"})

	require.Error(t, err)
}

func TestGitAdapterGenerateShouldReturnAnErrorWhenGivenAnInvalidOption(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	_, err = adapter.Generate([]string{"iaminvalid"})

	require.Error(t, err)
}

func TestGitAdapterGenerateShouldCreateAGitignoreFileWhenGivenASingleOption(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	contents, err := adapter.Generate([]string{"C"})

	require.NoError(t, err)
	require.Contains(t, contents, "*.o")
}

func TestGitAdapterGenerateShouldCreateAGitignoreFileWhenGivenMultipleOptions(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	contents, err := adapter.Generate([]string{"C", "Python"})

	require.NoError(t, err)
	require.Contains(t, contents, "### C ###")
	require.Contains(t, contents, "*.o")
	require.Contains(t, contents, "### Python ###")
	require.Contains(t, contents, "__pycache__/")
}

func TestGitAdapterGenerateShouldBeAbleToReadFromNestedDirectories(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	contents, err := adapter.Generate([]string{"Nikola"})

	require.NoError(t, err)
	require.Contains(t, contents, ".doit.db")
}

func TestGitAdapterGenerateShouldNotContainDuplicateHeaders(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	contents, err := adapter.Generate([]string{"Hugo"})

	require.NoError(t, err)
	require.Equal(t, 1, strings.Count(contents, "### Hugo ###"))
}

func TestGitAdapterUpdateShouldPUllTheRepository(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	repoPath := path.Join(testDir, "gitignore")
	require.DirExists(t, repoPath)
}

func TestGitAdapterUpdateShouldReturnAnErrorWhenTheRepoDoesNotExist(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: testDir,
		RepoURL:       "https://example.com/thisdoes/notexist.git",
	}

	err = adapter.Update()
	require.Error(t, err)
}

func TestGitAdapterUpdateShouldPullChangesToAnExistingRepo(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	err = adapter.Update()
	require.NoError(t, err)

	err = adapter.Update()
	require.NoError(t, err)

	repoPath := path.Join(testDir, "gitignore")
	require.DirExists(t, repoPath)
}

func TestGitAdapterUpdateShouldReturnAnErrorWhenThePathGivesAFileNotADir(t *testing.T) {
	t.Parallel()

	testDir, err := os.MkdirTemp("", "test-gitignore")
	if err != nil {
		t.Errorf("Unable to create temporary directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	adapter := &internal.GitAdapter{
		RepoDirectory: path.Join(testDir, "gitignore"),
		RepoURL:       internal.DefaultGitRepo,
	}

	_ = os.RemoveAll(testDir)

	file, err := os.Create(testDir)
	defer func() {
		_ = file.Close()
	}()
	require.NoError(t, err)

	err = adapter.Update()
	require.Error(t, err)
}
