package getters

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gotgenes/getignore/contentstructs"
	"github.com/gotgenes/getignore/errors"
	"github.com/gotgenes/getignore/testutils"
)

func TestGetIgnoreFilesForNameOnly(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]string{"Global/Vim"},
		[]contentstructs.NamedIgnoreContents{
			{
				"Global/Vim",
				".*.swp\nSession.vim\n",
			},
		},
		nil,
	)
}

func TestGetIgnoreFilesWithDefaultExtension(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]string{"Global/Vim.gitignore"},
		[]contentstructs.NamedIgnoreContents{
			{
				"Global/Vim.gitignore",
				".*.swp\nSession.vim\n",
			},
		},
		nil,
	)
}

func TestGetIgnoreFilesWithDifferentExtension(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]string{"Foo.bar"},
		[]contentstructs.NamedIgnoreContents{
			{
				"Foo.bar",
				"abc\nxyz\n",
			},
		},
		nil,
	)
}

func TestGetIgnoreFilesNotFound(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	expectedError := errors.FailedSources{
		errors.FailedSource{Source: testServer.URL + "/Nonexistent.gitignore", Err: fmt.Errorf("Got status code 404")},
	}
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]string{"Nonexistent"},
		nil,
		expectedError,
	)
}

type contentsAndError struct {
	Contents []contentstructs.NamedIgnoreContents
	Err      error
}

func TestGetIgnoreFilesContentsInExpectedOrder(t *testing.T) {
	names := []string{"Foo.bar", "Go.gitignore"}
	handler, responseGates := makeGatedHandler(names)
	testServer := httptest.NewServer(handler)
	defer testServer.Close()
	getter := HTTPGetter{
		testServer.URL,
		".gitignore",
		2,
	}
	expectedContents := []contentstructs.NamedIgnoreContents{
		{
			Name:     "Foo.bar",
			Contents: "abc\nxyz\n",
		},
		{
			Name:     "Go.gitignore",
			Contents: "*.o\n*.a\n*.so\n",
		},
	}
	resultChannel := make(chan contentsAndError)
	go func() {
		gotContents, err := getter.GetIgnoreFiles(names)
		resultChannel <- contentsAndError{gotContents, err}
	}()
	responseGates["Go.gitignore"] <- true
	responseGates["Foo.bar"] <- true
	result := <-resultChannel
	testutils.AssertDeepEqual(t, result.Contents, expectedContents)
	testutils.AssertDeepEqual(t, result.Err, nil)
}

func assertGetIgnoreFilesReturnsExpectedContents(t *testing.T, testServer *httptest.Server, names []string, expectedContents []contentstructs.NamedIgnoreContents, expectedError error) {
	getter := HTTPGetter{
		testServer.URL,
		".gitignore",
		1,
	}
	gotContents, err := getter.GetIgnoreFiles(names)
	testutils.AssertDeepEqual(t, gotContents, expectedContents)
	testutils.AssertDeepEqual(t, err, expectedError)
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

func makeGatedHandler(names []string) (handler http.Handler, responseGates map[string]chan bool) {
	responseGates = make(map[string]chan bool)
	for _, name := range names {
		responseGates[name] = make(chan bool)
	}
	hf := func(w http.ResponseWriter, r *http.Request) {
		gatedHandlerFunc(w, r, responseGates)
	}
	handler = http.HandlerFunc(hf)
	return
}

func gatedHandlerFunc(w http.ResponseWriter, r *http.Request, responseGates map[string]chan bool) {
	name := r.URL.Path[1:]
	gate := responseGates[name]
	<-gate
	contents, ok := pathsToContents[name]
	if ok {
		fmt.Fprint(w, contents)
	} else {
		w.WriteHeader(404)
	}
}
