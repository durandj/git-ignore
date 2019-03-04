package gitignore_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/durandj/git-ignore/pkg/gitignore"
)

var _ = Describe("Client", func() {
	var primaryAdapter fakeAdapter
	var secondaryAdapter fakeAdapter
	var client gitignore.Client

	BeforeEach(func() {
		primaryAdapter = newFakeAdapter()
		secondaryAdapter = newFakeAdapter()

		client = gitignore.Client{
			Adapters: []gitignore.Adapter{
				&primaryAdapter,
				&secondaryAdapter,
			},
		}
	})

	Describe("List", func() {
		Context("with error in one adapter", func() {
			It("should fallback to the next adapter", func() {
				expectedOptions := []string{"c"}

				primaryAdapter.addListReturn(nil, fmt.Errorf("Primary error"))
				secondaryAdapter.addListReturn(expectedOptions, nil)

				options, err := client.List()

				Expect(err).To(BeNil())

				Expect(options).To(ConsistOf(expectedOptions))
			})
		})

		Context("with error in all adapters", func() {
			It("should return an error", func() {
				expectedErr := fmt.Errorf("Test error")
				primaryAdapter.addListReturn(nil, expectedErr)
				secondaryAdapter.addListReturn(nil, expectedErr)

				_, err := client.List()

				Expect(err).ToNot(BeNil())
			})
		})

		It("should retrieve a list of options", func() {
			expectedOptions := []string{"c", "c++"}
			primaryAdapter.addListReturn(expectedOptions, nil)

			options, err := client.List()

			Expect(err).To(BeNil())

			Expect(options).ToNot(BeNil())
			Expect(options).ToNot(BeEmpty())
			Expect(options).To(ConsistOf(expectedOptions))
		})
	})

	Describe("Generate", func() {
		Context("with no options", func() {
			It("should return an error", func() {
				_, err := client.Generate(nil)

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with an invalid option", func() {
			It("should return an error", func() {
				primaryAdapter.addListReturn([]string{"c"}, nil)

				_, err := client.Generate([]string{"doesnotexist"})

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with a single option", func() {
			It("should generate a gitignore file", func() {
				primaryAdapter.addListReturn([]string{"c", "c++"}, nil)
				primaryAdapter.addGenerateReturn("### C ###", nil)

				file, err := client.Generate([]string{"c"})

				Expect(err).To(BeNil())

				Expect(file).To(ContainSubstring("### C ###"))

				generateCalls := primaryAdapter.getGenerateCalls()
				Expect(generateCalls).To(HaveLen(1))
				Expect(generateCalls[0].options).To(ConsistOf("c"))
			})
		})

		Context("with multiple options", func() {
			It("should generate a gitignore file", func() {
				primaryAdapter.addListReturn([]string{"c", "c++"}, nil)
				primaryAdapter.addGenerateReturn("### C ###\n\n### C++ ###", nil)

				file, err := client.Generate([]string{"c", "c++"})

				Expect(err).To(BeNil())

				Expect(file).To(ContainSubstring("### C ###"))
				Expect(file).To(ContainSubstring("### C++ ###"))

				generateCalls := primaryAdapter.getGenerateCalls()
				Expect(generateCalls).To(HaveLen(1))
				Expect(generateCalls[0].options).To(ConsistOf("c", "c++"))
			})
		})

		Context("with an error in one adapter", func() {
			It("should fallback to the next adapter", func() {
				primaryAdapter.addListReturn(nil, fmt.Errorf("Test error"))
				secondaryAdapter.addListReturn([]string{"c"}, nil)
				secondaryAdapter.addGenerateReturn("### C ###", nil)

				file, err := client.Generate([]string{"c"})

				Expect(err).To(BeNil())

				Expect(file).To(ContainSubstring("### C ###"))

				primaryListCalls := primaryAdapter.getListCalls()
				primaryGenerateCalls := primaryAdapter.getGenerateCalls()
				secondaryListCalls := secondaryAdapter.getListCalls()
				secondaryGenerateCalls := secondaryAdapter.getGenerateCalls()
				Expect(primaryListCalls).To(HaveLen(1))
				Expect(primaryGenerateCalls).To(BeEmpty())
				Expect(secondaryListCalls).To(HaveLen(1))
				Expect(secondaryGenerateCalls).To(HaveLen(1))
				Expect(secondaryGenerateCalls[0].options).To(ConsistOf("c"))
			})
		})

		Context("with an error in all adapters", func() {
			It("should return an error", func() {
				primaryAdapter.addListReturn(nil, fmt.Errorf("Test error"))
				secondaryAdapter.addListReturn(nil, fmt.Errorf("Test error"))

				_, err := client.Generate([]string{"c"})

				Expect(err).ToNot(BeNil())
			})
		})
	})

	Describe("Update", func() {
		It("should update all supporting adapters", func() {
			primaryMappings := make(map[string]string)
			secondaryMappings := make(map[string]string)

			primaryMappings["c"] = "### C ###"
			secondaryMappings["c++"] = "### C++ ###"

			primaryAdapter.addSourceReturn(primaryMappings, nil)
			secondaryAdapter.addSourceReturn(secondaryMappings, nil)

			primaryAdapter.addCacheReturn(nil)
			secondaryAdapter.addCacheReturn(nil)

			err := client.Update()

			Expect(err).To(BeNil())

			primarySourceCalls := primaryAdapter.getSourceCalls()
			primaryCacheCalls := primaryAdapter.getCacheCalls()
			secondarySourceCalls := secondaryAdapter.getSourceCalls()
			secondaryCacheCalls := secondaryAdapter.getCacheCalls()

			Expect(primarySourceCalls).To(HaveLen(1))
			Expect(secondarySourceCalls).To(HaveLen(1))
			Expect(primaryCacheCalls).To(HaveLen(1))
			Expect(primaryCacheCalls[0].mappings).To(HaveLen(2))
			Expect(primaryCacheCalls[0].mappings).To(HaveKey("c"))
			Expect(primaryCacheCalls[0].mappings["c"]).To(Equal("### C ###"))
			Expect(primaryCacheCalls[0].mappings).To(HaveKey("c++"))
			Expect(primaryCacheCalls[0].mappings["c++"]).To(Equal("### C++ ###"))
			Expect(secondaryCacheCalls).To(HaveLen(1))
			Expect(secondaryCacheCalls[0].mappings).To(HaveLen(2))
			Expect(secondaryCacheCalls[0].mappings).To(HaveKey("c"))
			Expect(secondaryCacheCalls[0].mappings["c"]).To(Equal("### C ###"))
			Expect(secondaryCacheCalls[0].mappings).To(HaveKey("c++"))
			Expect(secondaryCacheCalls[0].mappings["c++"]).To(Equal("### C++ ###"))
		})

		Context("with an error in one or more adapters during sourcing", func() {
			It("should return an error", func() {
				secondaryMappings := make(map[string]string)
				secondaryMappings["c++"] = "### C++ ###"

				primaryAdapter.addSourceReturn(nil, fmt.Errorf("Test error"))
				secondaryAdapter.addSourceReturn(secondaryMappings, nil)

				primaryAdapter.addCacheReturn(nil)
				secondaryAdapter.addCacheReturn(nil)

				err := client.Update()

				Expect(err).ToNot(BeNil())

				Expect(primaryAdapter.getCacheCalls()).To(HaveLen(0))
				Expect(secondaryAdapter.getCacheCalls()).To(HaveLen(0))
			})
		})

		Context("with an error in one or more adapters during caching", func() {
			It("should continue caching with working adapters", func() {
				primaryMappings := make(map[string]string)
				secondaryMappings := make(map[string]string)

				primaryMappings["c"] = "### C ###"
				secondaryMappings["c++"] = "### C++ ###"

				primaryAdapter.addSourceReturn(primaryMappings, nil)
				secondaryAdapter.addSourceReturn(secondaryMappings, nil)

				primaryAdapter.addCacheReturn(fmt.Errorf("Test error"))
				secondaryAdapter.addCacheReturn(nil)

				err := client.Update()

				Expect(err).To(BeNil())

				primaryCacheCalls := primaryAdapter.getCacheCalls()
				secondaryCacheCalls := secondaryAdapter.getCacheCalls()

				Expect(primaryCacheCalls).To(HaveLen(1))
				Expect(secondaryCacheCalls).To(HaveLen(1))
			})
		})
	})
})
