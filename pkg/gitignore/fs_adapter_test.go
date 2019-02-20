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

func createTestFile(testDir string, name string, contents string) {
	file, err := os.Create(path.Join(testDir, name))
	if err != nil {
		fmt.Printf("Unable to create test file: %s\n", err)
		os.Exit(1)
	}

	defer file.Close()

	file.WriteString(contents)
}

var _ = Describe("FSAdapter", func() {
	var testDir string
	var adapter *gitignore.FSAdapter

	BeforeEach(func() {
		var err error
		testDir, err = ioutil.TempDir("", "test-gitignore-")
		if err != nil {
			fmt.Printf("Unable to create temporary directory: %s\n", err)
			os.Exit(1)
		}

		adapter, err = gitignore.NewFSAdapter(testDir)
		if err != nil {
			fmt.Printf("Unable to create file system adapter: %s\n", err)
			os.Exit(1)
		}
	})

	AfterEach(func() {
		os.RemoveAll(testDir)
	})

	Describe("List", func() {
		Context("with no existing files", func() {
			It("should return an error", func() {
				os.RemoveAll(testDir)

				_, err := adapter.List()

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with existing files", func() {
			It("should return options from files", func() {
				createTestFile(testDir, "c", "### C ###")
				createTestFile(testDir, "c++", "### C++ ###")

				options, err := adapter.List()

				Expect(err).To(BeNil())

				Expect(options).To(ConsistOf("c", "c++"))
			})
		})
	})

	Describe("Generate", func() {
		Context("with no existing files", func() {
			It("should return an error", func() {
				_, err := adapter.Generate([]string{"c"})

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with partial file cache", func() {
			It("should return an error", func() {
				createTestFile(testDir, "c", "### C ###")

				_, err := adapter.Generate([]string{"c", "c++"})

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with full existing file cache", func() {
			It("should return file contents", func() {
				createTestFile(testDir, "c", "### C ###")
				createTestFile(testDir, "c++", "### C++ ###")

				contents, err := adapter.Generate([]string{"c", "c++"})

				Expect(err).To(BeNil())

				Expect(contents).To(ContainSubstring("### C ###"))
				Expect(contents).To(ContainSubstring("### C++ ###"))
			})
		})
	})
})
