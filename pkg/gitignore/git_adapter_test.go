package gitignore_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/durandj/git-ignore/pkg/gitignore"
)

func directoryExists(filePath string) bool {
	fileStat, err := os.Stat(filePath)

	return err == nil && fileStat.IsDir()
}

var _ = Describe("GitAdapter", func() {
	var testDir string
	var adapter *gitignore.GitAdapter

	BeforeEach(func() {
		var err error
		testDir, err = ioutil.TempDir("", "test-gitignore")
		if err != nil {
			fmt.Printf("Unable to create temporary directory: %s\n", err)
			os.Exit(1)
		}

		adapter = &gitignore.GitAdapter{
			RepoDirectory: path.Join(testDir, "gitignore"),
			RepoURL:       gitignore.DefaultGitRepo,
		}
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	Describe("List", func() {
		It("should retrieve a list of options", func() {
			err := adapter.Update()
			Expect(err).To(BeNil())

			options, err := adapter.List()
			Expect(err).To(BeNil())

			Expect(len(options)).To(BeNumerically(">", 0))
			Expect(options).To(ContainElement("C"))
			Expect(options).To(ContainElement("C++"))
		})

		It("should return an error if the repository doesn't exist", func() {
			_, err := adapter.List()

			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Generate", func() {
		BeforeEach(func() {
			if err := adapter.Update(); err != nil {
				panic(err)
			}
		})

		It("should return an error when no options are given", func() {
			_, err := adapter.Generate([]string{})

			Expect(err).ToNot(BeNil())
		})

		It("should return an error when the repository doesn't exist", func() {
			os.RemoveAll(testDir)

			_, err := adapter.Generate([]string{"C"})

			Expect(err).ToNot(BeNil())
		})

		It("should return an error when given an invalid option", func() {
			_, err := adapter.Generate([]string{"iaminvalid"})

			Expect(err).ToNot(BeNil())
		})

		It("should generate a gitignore file when given a single option", func() {
			contents, err := adapter.Generate([]string{"C"})

			Expect(err).To(BeNil())

			Expect(contents).To(ContainSubstring("*.o"))
		})

		It("should generate a gitignore file when given multiple options", func() {
			contents, err := adapter.Generate([]string{"C", "Python"})

			Expect(err).To(BeNil())

			Expect(contents).To(ContainSubstring("### C ###"))
			Expect(contents).To(ContainSubstring("*.o"))
			Expect(contents).To(ContainSubstring("### Python ###"))
			Expect(contents).To(ContainSubstring("__pycache__/"))
		})

		It("should be able to read from nested directories", func() {
			contents, err := adapter.Generate([]string{"Nikola"})

			Expect(err).To(BeNil())

			Expect(contents).To(ContainSubstring(".doit.db"))
		})
	})

	Describe("Update", func() {
		It("should pull the repository", func() {
			err := adapter.Update()

			Expect(err).To(BeNil())

			repoPath := path.Join(testDir, "gitignore")
			Expect(directoryExists(repoPath)).To(BeTrue())
		})

		It("should return an error when the repository doesn't exist", func() {
			adapter = &gitignore.GitAdapter{
				RepoDirectory: testDir,
				RepoURL:       "https://example.com/thisdoes/notexist.git",
			}

			err := adapter.Update()

			Expect(err).ToNot(BeNil())
		})

		It("should pull changes to an existing repository", func() {
			err := adapter.Update()
			Expect(err).To(BeNil())

			err = adapter.Update()
			Expect(err).To(BeNil())

			repoPath := path.Join(testDir, "gitignore")
			Expect(directoryExists(repoPath)).To(BeTrue())
		})

		It("should return an error when the path points at a file instead of a directory", func() {
			os.RemoveAll(testDir)

			file, err := os.Create(testDir)
			Expect(err).To(BeNil())
			file.Close()

			err = adapter.Update()
			Expect(err).ToNot(BeNil())
		})
	})
})
