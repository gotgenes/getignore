package getters

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gotgenes/getignore/contentstructs"
	"github.com/gotgenes/getignore/testutils"
)

func TestGetIgnoreFilesForNameOnly(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]string{"Global/Vim"},
		[]contentstructs.RetrievedContents{
			{
				"Global/Vim",
				testServer.URL + "/Global/Vim.gitignore",
				".*.swp\nSession.vim\n",
				nil,
			},
		},
	)
}

func TestGetIgnoreFilesWithDefaultExtension(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]string{"Global/Vim.gitignore"},
		[]contentstructs.RetrievedContents{
			{
				"Global/Vim.gitignore",
				testServer.URL + "/Global/Vim.gitignore",
				".*.swp\nSession.vim\n",
				nil,
			},
		},
	)
}

func TestGetIgnoreFilesWithDifferentExtension(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]string{"Foo.bar"},
		[]contentstructs.RetrievedContents{
			{
				"Foo.bar",
				testServer.URL + "/Foo.bar",
				"abc\nxyz\n",
				nil,
			},
		},
	)
}

func TestGetIgnoreFilesNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]string{"Nonexistent"},
		[]contentstructs.RetrievedContents{
			{
				"Nonexistent",
				testServer.URL + "/Nonexistent.gitignore",
				"",
				fmt.Errorf("Got status code 404"),
			},
		},
	)
}

func assertGetIgnoreFilesReturnsExpectedContents(t *testing.T, testServer *httptest.Server, names []string, allExpectedContents []contentstructs.RetrievedContents) {
	getter := HTTPGetter{
		testServer.URL,
		".gitignore",
		1,
	}
	contentsChannel := make(chan contentstructs.RetrievedContents)
	var requestsPending sync.WaitGroup
	getter.GetIgnoreFiles(names, contentsChannel, &requestsPending)
	for _, expectedContents := range allExpectedContents {
		gotContents := <-contentsChannel
		testutils.AssertDeepEqual(t, gotContents, expectedContents)
	}
}

var pathsToContents = map[string]string{
	"Global/Vim.gitignore": ".*.swp\nSession.vim\n",
	"Go.gitignore":         "*.o\n*.a\n*.so\n",
	"Foo.bar":              "abc\nxyz\n",
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	contents, ok := pathsToContents[r.URL.Path[1:]]
	if ok {
		fmt.Fprint(w, contents)
	} else {
		w.WriteHeader(404)
	}
}
