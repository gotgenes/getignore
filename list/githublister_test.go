package list_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gotgenes/getignore/identifiers"
	"github.com/gotgenes/getignore/list"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("GitHubLister", func() {
	var (
		ctx    context.Context
		server *ghttp.Server
		lister list.GitHubLister
	)

	BeforeEach(func() {
		ctx = context.Background()
		server = ghttp.NewServer()
		lister, _ = list.NewGitHubLister(list.WithBaseURL(server.URL()))
	})

	AfterEach(func() {
		server.Close()
	})

	Context("basic functionality", func() {
		expectedUserAgent := []string{fmt.Sprintf("getignore/%s", identifiers.Version)}
		BeforeEach(func() {
			responseBody := `{
	  "name": "master",
	  "commit": {
		"sha": "b0012e4930d0a8c350254a3caeedf7441ea286a3",
		"commit": {
		  "tree": {
			"sha": "5adf061bdde4dd26889be1e74028b2f54aabc346"
		  }
		}
	  }
	}`
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v3/repos/github/gitignore/branches/master"),
					ghttp.VerifyHeader(http.Header{
						"User-Agent": expectedUserAgent,
					}),
					ghttp.VerifyHeader(http.Header{
						"Accept": []string{"application/vnd.github.v3+json"},
					}),
					ghttp.RespondWith(http.StatusOK, responseBody),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v3/repos/github/gitignore/git/trees/5adf061bdde4dd26889be1e74028b2f54aabc346"),
					ghttp.VerifyHeader(http.Header{
						"User-Agent": expectedUserAgent,
					}),
					ghttp.VerifyHeader(http.Header{
						"Accept": []string{"application/vnd.github.v3+json"},
					}),
				),
			)

		})

		It("should send requests with the expected headers", func() {

			lister.List(ctx, "")

			Expect(server.ReceivedRequests()).Should(HaveLen(2))
		})
	})

	Context("happy path", func() {

		BeforeEach(func() {
			branchesResponseBody := `{
	  "name": "master",
	  "commit": {
		"sha": "b0012e4930d0a8c350254a3caeedf7441ea286a3",
		"commit": {
		  "tree": {
			"sha": "5adf061bdde4dd26889be1e74028b2f54aabc346"
		  }
		}
	  }
	}`
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v3/repos/github/gitignore/branches/master"),
					ghttp.RespondWith(http.StatusOK, branchesResponseBody),
				),
			)
		})

		When("the tree response is empty", func() {

			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v3/repos/github/gitignore/git/trees/5adf061bdde4dd26889be1e74028b2f54aabc346"),
						ghttp.RespondWith(http.StatusOK, "{}"),
					),
				)
			})

			It("should return an empty slice", func() {
				ignoreFiles, err := lister.List(ctx, "")
				Expect(ignoreFiles).Should(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
