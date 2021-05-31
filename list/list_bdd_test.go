package list_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gotgenes/getignore/identifiers"
	"github.com/gotgenes/getignore/list"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lister", func() {
	var (
		server *httptest.Server
		client *http.Client
	)

	It("should send a request with the expected headers", func() {
		var header http.Header
		server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header = r.Header
			w.Write([]byte("Yay!"))
		}))
		client = server.Client()
		lister := list.NewGitHubLister(*client, server.URL)
		lister.List("")

		expectedAccept := []string{"application/vnd.github.v3+json"}
		expectedUserAgent := []string{fmt.Sprintf("getignore/%s", identifiers.Version)}
		Expect(header.Values("Accept")).To(Equal(expectedAccept))
		Expect(header.Values("User-Agent")).To(Equal(expectedUserAgent))
	})

})
