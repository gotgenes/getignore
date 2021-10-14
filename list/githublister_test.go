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

			lister.List(ctx)

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
				ignoreFiles, _ := lister.List(ctx)
				Expect(ignoreFiles).Should(BeNil())
			})

			It("should not have an error", func() {
				_, err := lister.List(ctx)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("the response has gitignore files", func() {
			BeforeEach(func() {
				responseBody := `{
  "sha": "5adf061bdde4dd26889be1e74028b2f54aabc346",
  "url": "https://api.github.com/repos/github/gitignore/git/trees/5adf061bdde4dd26889be1e74028b2f54aabc346",
  "tree": [
    {
      "path": "Actionscript.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "5d947ca8879f8a9072fe485c566204e3c2929e80",
      "size": 350,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/5d947ca8879f8a9072fe485c566204e3c2929e80"
    },
    {
      "path": "Global/Anjuta.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "20dd42c53e6f0df8233fee457b664d443ee729f4",
      "size": 78,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/20dd42c53e6f0df8233fee457b664d443ee729f4"
    },
    {
      "path": "community/AWS/SAM.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "dc9d020aee1ebc1a23c02d80a1c33c0cb35ebaeb",
      "size": 167,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/dc9d020aee1ebc1a23c02d80a1c33c0cb35ebaeb"
    }
  ],
  "truncated": false
}
				}`
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v3/repos/github/gitignore/git/trees/5adf061bdde4dd26889be1e74028b2f54aabc346"),
						ghttp.RespondWith(http.StatusOK, responseBody),
					),
				)
			})

			It("should return a list of gitignore files", func() {
				ignoreFiles, _ := lister.List(ctx)
				Expect(ignoreFiles).Should(Equal(
					[]string{"Actionscript.gitignore", "Global/Anjuta.gitignore", "community/AWS/SAM.gitignore"}))
			})

			It("should not have an error", func() {
				_, err := lister.List(ctx)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("the response has additional files", func() {
			BeforeEach(func() {
				responseBody := `{
  "sha": "5adf061bdde4dd26889be1e74028b2f54aabc346",
  "url": "https://api.github.com/repos/github/gitignore/git/trees/5adf061bdde4dd26889be1e74028b2f54aabc346",
  "tree": [
    {
      "path": ".github/PULL_REQUEST_TEMPLATE.md",
      "mode": "100644",
      "type": "blob",
      "sha": "247a5b56e890c2ab29eb337f26aa623deb2feefc",
      "size": 199,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/247a5b56e890c2ab29eb337f26aa623deb2feefc"
    },
    {
      "path": ".travis.yml",
      "mode": "100644",
      "type": "blob",
      "sha": "4009e0bc8b07582c19fa761810c9f3741ab76597",
      "size": 103,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/4009e0bc8b07582c19fa761810c9f3741ab76597"
    },
    {
      "path": "Actionscript.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "5d947ca8879f8a9072fe485c566204e3c2929e80",
      "size": 350,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/5d947ca8879f8a9072fe485c566204e3c2929e80"
    },
    {
      "path": "Global/Anjuta.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "20dd42c53e6f0df8233fee457b664d443ee729f4",
      "size": 78,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/20dd42c53e6f0df8233fee457b664d443ee729f4"
    },
    {
      "path": "community/AWS/SAM.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "dc9d020aee1ebc1a23c02d80a1c33c0cb35ebaeb",
      "size": 167,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/dc9d020aee1ebc1a23c02d80a1c33c0cb35ebaeb"
    }
  ],
  "truncated": false
}
				}`
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v3/repos/github/gitignore/git/trees/5adf061bdde4dd26889be1e74028b2f54aabc346"),
						ghttp.RespondWith(http.StatusOK, responseBody),
					),
				)
			})

			It("should filter files with .gitignore suffix", func() {
				ignoreFiles, _ := lister.List(ctx)
				Expect(ignoreFiles).Should(Equal(
					[]string{"Actionscript.gitignore", "Global/Anjuta.gitignore", "community/AWS/SAM.gitignore"}))
			})

			It("should not have an error", func() {
				_, err := lister.List(ctx)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("the response has directories", func() {
			BeforeEach(func() {
				responseBody := `{
  "sha": "5adf061bdde4dd26889be1e74028b2f54aabc346",
  "url": "https://api.github.com/repos/github/gitignore/git/trees/5adf061bdde4dd26889be1e74028b2f54aabc346",
  "tree": [
    {
      "path": ".github",
      "mode": "040000",
      "type": "tree",
      "sha": "45f58ef9211cc06f3ef86585c7ecb1b3d52fd4f9",
      "url": "https://api.github.com/repos/github/gitignore/git/trees/45f58ef9211cc06f3ef86585c7ecb1b3d52fd4f9"
    },
    {
      "path": "Actionscript.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "5d947ca8879f8a9072fe485c566204e3c2929e80",
      "size": 350,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/5d947ca8879f8a9072fe485c566204e3c2929e80"
    },
    {
      "path": "foo.gitignore",
      "mode": "040000",
      "type": "tree",
      "sha": "a1f9ba2be789d9d7a3559967c42f22cbea9bf8dc",
      "url": "https://api.github.com/repos/github/gitignore/git/trees/a1f9ba2be789d9d7a3559967c42f22cbea9bf8dc"
    },
    {
      "path": "Global",
      "mode": "040000",
      "type": "tree",
      "sha": "5fb11fe033ab0f8a86b7b5aa8e4f13f9d5d3f7ca",
      "url": "https://api.github.com/repos/github/gitignore/git/trees/5fb11fe033ab0f8a86b7b5aa8e4f13f9d5d3f7ca"
    },
    {
      "path": "Global/Anjuta.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "20dd42c53e6f0df8233fee457b664d443ee729f4",
      "size": 78,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/20dd42c53e6f0df8233fee457b664d443ee729f4"
    },
    {
      "path": "community/AWS/SAM.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "dc9d020aee1ebc1a23c02d80a1c33c0cb35ebaeb",
      "size": 167,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/dc9d020aee1ebc1a23c02d80a1c33c0cb35ebaeb"
    }
  ],
  "truncated": false
}
				}`
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/api/v3/repos/github/gitignore/git/trees/5adf061bdde4dd26889be1e74028b2f54aabc346"),
						ghttp.RespondWith(http.StatusOK, responseBody),
					),
				)
			})

			It("should return only files", func() {
				ignoreFiles, _ := lister.List(ctx)
				Expect(ignoreFiles).Should(Equal(
					[]string{"Actionscript.gitignore", "Global/Anjuta.gitignore", "community/AWS/SAM.gitignore"}))
			})

			It("should not have an error", func() {
				_, err := lister.List(ctx)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
