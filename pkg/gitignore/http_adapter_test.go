package gitignore_test

import (
	"fmt"
	"net/http"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/durandj/git-ignore/pkg/gitignore"
)

type generateResponse struct {
	name    string
	content string
}

func createGenerateResponse(responses []generateResponse) string {
	options := make([]string, len(responses))
	for index, response := range responses {
		options[index] = response.name
	}
	optionsStr := strings.Join(options, ",")

	var builder strings.Builder

	builder.WriteString(
		fmt.Sprintf("# Created by https://www.gitignore.io/api/%s\n", optionsStr),
	)
	builder.WriteString(
		fmt.Sprintf("# Edit at https://www.gitignore.io/?templates=%s\n", optionsStr),
	)

	for _, response := range responses {
		builder.WriteString(response.content)
		builder.WriteString("\n")
	}

	builder.WriteString(
		fmt.Sprintf("# End of https://www.gitignore.io/api/%s\n", optionsStr),
	)

	return builder.String()
}

var _ = Describe("HTTPAdapter", func() {
	var server *ghttp.Server
	var adapter *gitignore.HTTPAdapter

	BeforeEach(func() {
		server = ghttp.NewServer()
		adapter = gitignore.NewHTTPAdapter(server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("List", func() {
		It("should retrieve a list of options", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/list"),
					ghttp.RespondWith(http.StatusOK, "c,c++\npython"),
				),
			)

			options, err := adapter.List()

			Expect(err).To(BeNil())

			Expect(options).ToNot(BeNil())
			Expect(options).ToNot(BeEmpty())
			Expect(options).To(ContainElement("c"))
		})

		It("should return an error if the request fails", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/list"),
					ghttp.RespondWith(http.StatusGatewayTimeout, "Unavailable"),
				),
			)

			_, err := adapter.List()

			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Generate", func() {
		Context("with no options", func() {
			It("should return an error", func() {
				_, err := adapter.Generate(nil)

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with an failed request", func() {
			It("should return an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/list"),
						ghttp.RespondWith(http.StatusOK, "c,c++"),
					),

					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/c"),
						ghttp.RespondWith(http.StatusGatewayTimeout, "Unavailable"),
					),
				)

				_, err := adapter.Generate([]string{"c"})

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with an invalid option", func() {
			It("should return an error", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/list"),
						ghttp.RespondWith(http.StatusOK, "c,c++"),
					),
				)

				_, err := adapter.Generate([]string{"doesnotexist"})

				Expect(err).ToNot(BeNil())
			})
		})

		Context("with a single option", func() {
			It("should generate a gitignore file", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/list"),
						ghttp.RespondWith(http.StatusOK, "c,c++"),
					),

					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/c"),
						ghttp.RespondWith(http.StatusOK, "### C ###"),
					),
				)

				file, err := adapter.Generate([]string{"c"})

				Expect(err).To(BeNil())

				Expect(file).To(ContainSubstring("### C ###"))
			})
		})

		Context("with multiple options", func() {
			It("should generate a gitignore file", func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/list"),
						ghttp.RespondWith(http.StatusOK, "c,c++"),
					),

					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/c,c++"),
						ghttp.RespondWith(
							http.StatusOK,
							createGenerateResponse([]generateResponse{
								{name: "c", content: "### C ###"},
								{name: "c++", content: "### C++ ###"},
							}),
						),
					),
				)

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

	Describe("Source", func() {
		It("should generate a mapping of different options to contents", func() {
			server.RouteToHandler(
				"GET",
				"/list",
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/list"),
					ghttp.RespondWith(http.StatusOK, "c,c++"),
				),
			)
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/c"),
					ghttp.RespondWith(http.StatusOK, "### C ###"),
				),

				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/c++"),
					ghttp.RespondWith(http.StatusOK, "### C++ ###"),
				),
			)

			mappings, err := adapter.Source()

			Expect(err).To(BeNil())

			Expect(mappings).ToNot(BeEmpty())

			Expect(mappings).To(HaveKey("c"))
			Expect(mappings["c"]).To(ContainSubstring("### C ###"))

			Expect(mappings).To(HaveKey("c++"))
			Expect(mappings["c++"]).To(ContainSubstring("### C++ ###"))
		})
	})
})
