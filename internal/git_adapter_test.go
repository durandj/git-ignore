package internal_test

import (
	"fmt"
	"os"
	"path"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/durandj/git-ignore/internal"
)

func directoryExists(filePath string) bool {
	fileStat, err := os.Stat(filePath)

	return err == nil && fileStat.IsDir()
}

var _ = Describe("GitAdapter", func() {
	var testDir string
	var adapter *internal.GitAdapter

	BeforeEach(func() {
		var err error
		testDir, err = os.MkdirTemp("", "test-gitignore")
		if err != nil {
			Fail(fmt.Sprintf("Unable to create temporary directory: %v", err))
		}

		adapter = &internal.GitAdapter{
			RepoDirectory: path.Join(testDir, "gitignore"),
			RepoURL:       internal.DefaultGitRepo,
		}
	})

	AfterEach(func() {
		_ = os.RemoveAll(testDir)
	})

	Describe("List", func() {
		It("should retrieve a list of options", func() {
			err := adapter.Update()
			Expect(err).ToNot(HaveOccurred())

			options, err := adapter.List()
			Expect(err).ToNot(HaveOccurred())

			Expect(len(options)).NotTo(BeEmpty())
			Expect(options).To(ContainElement("C"))
			Expect(options).To(ContainElement("C++"))
		})

		It("should return an error if the repository doesn't exist", func() {
			_, err := adapter.List()

			Expect(err).To(HaveOccurred())
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

			Expect(err).To(HaveOccurred())
		})

		It("should return an error when the repository doesn't exist", func() {
			_ = os.RemoveAll(testDir)

			_, err := adapter.Generate([]string{"C"})

			Expect(err).To(HaveOccurred())
		})

		It("should return an error when given an invalid option", func() {
			_, err := adapter.Generate([]string{"iaminvalid"})

			Expect(err).To(HaveOccurred())
		})

		It("should generate a gitignore file when given a single option", func() {
			contents, err := adapter.Generate([]string{"C"})

			Expect(err).ToNot(HaveOccurred())

			Expect(contents).To(ContainSubstring("*.o"))
		})

		It("should generate a gitignore file when given multiple options", func() {
			contents, err := adapter.Generate([]string{"C", "Python"})

			Expect(err).ToNot(HaveOccurred())

			Expect(contents).To(ContainSubstring("### C ###"))
			Expect(contents).To(ContainSubstring("*.o"))
			Expect(contents).To(ContainSubstring("### Python ###"))
			Expect(contents).To(ContainSubstring("__pycache__/"))
		})

		It("should be able to read from nested directories", func() {
			contents, err := adapter.Generate([]string{"Nikola"})

			Expect(err).ToNot(HaveOccurred())

			Expect(contents).To(ContainSubstring(".doit.db"))
		})

		It("should not contain duplicate headers", func() {
			contents, err := adapter.Generate([]string{"Hugo"})

			Expect(err).ToNot(HaveOccurred())
			Expect(strings.Count(contents, "### Hugo ###")).To(Equal(1))
		})
	})

	Describe("Update", func() {
		It("should pull the repository", func() {
			err := adapter.Update()

			Expect(err).ToNot(HaveOccurred())

			repoPath := path.Join(testDir, "gitignore")
			Expect(directoryExists(repoPath)).To(BeTrue())
		})

		It("should return an error when the repository doesn't exist", func() {
			adapter = &internal.GitAdapter{
				RepoDirectory: testDir,
				RepoURL:       "https://example.com/thisdoes/notexist.git",
			}

			err := adapter.Update()

			Expect(err).To(HaveOccurred())
		})

		It("should pull changes to an existing repository", func() {
			err := adapter.Update()
			Expect(err).ToNot(HaveOccurred())

			err = adapter.Update()
			Expect(err).ToNot(HaveOccurred())

			repoPath := path.Join(testDir, "gitignore")
			Expect(directoryExists(repoPath)).To(BeTrue())
		})

		It("should return an error when the path points at a file instead of a directory", func() {
			_ = os.RemoveAll(testDir)

			file, err := os.Create(testDir)
			defer func() {
				_ = file.Close()
			}()
			Expect(err).ToNot(HaveOccurred())

			err = adapter.Update()
			Expect(err).To(HaveOccurred())
		})
	})
})
