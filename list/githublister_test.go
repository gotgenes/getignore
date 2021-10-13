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

	It("should send requests with the expected headers", func() {
		expectedUserAgent := []string{fmt.Sprintf("getignore/%s", identifiers.Version)}
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

		lister.List(ctx, "")

		Expect(server.ReceivedRequests()).Should(HaveLen(2))
	})
})
