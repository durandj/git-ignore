package gitignore_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/durandj/git-ignore/pkg/gitignore"
)

var _ = Describe("HTTPAdapter", func() {
	var adapter gitignore.Adapter

	BeforeEach(func() {
		adapter = &gitignore.HTTPAdapter{}
	})

	Describe("List", func() {
		It("should retrieve a list of options", func() {
			options, err := adapter.List()

			Expect(err).To(BeNil())

			Expect(options).ToNot(BeNil())
			Expect(options).ToNot(BeEmpty())
			Expect(options).To(ContainElement("c"))
		})
	})

	Describe("Generate", func() {
		Context("with no options", func() {
			It("should return an error", func() {
				_, err := adapter.Generate(nil)

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with an invalid option", func() {
			It("should return an error", func() {
				_, err := adapter.Generate([]string{"doesnotexist"})

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with a single option", func() {
			It("should generate a gitignore file", func() {
				file, err := adapter.Generate([]string{"c"})

				Expect(err).To(BeNil())

				Expect(file).To(ContainSubstring("### C ###"))
			})
		})

		Context("with multiple options", func() {
			It("should generate a gitignore file", func() {
				file, err := adapter.Generate([]string{"c", "c++"})

				Expect(err).To(BeNil())

				Expect(file).To(ContainSubstring("### C ###"))
				Expect(file).To(ContainSubstring("### C++ ###"))
				Expect(file).ToNot(ContainSubstring("# Created by https://www.gitignore.io/api/c,c++"))
				Expect(file).ToNot(ContainSubstring("# Edit at https://www.gitignore.io/?templates=c,c++"))
				Expect(file).ToNot(ContainSubstring("# End of https://www.gitignore.io/api/c,c++"))
			})
		})
	})
})
