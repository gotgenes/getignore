package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"
)

var errorTemplate = "Got %q, expected %q"

func TestParseNamesFile(t *testing.T) {
	namesFile := bytes.NewBufferString("Global/Vim\nPython\n")
	names := ParseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf(errorTemplate, names, expectedNames)
	}
}

func TestParseNamesFileIgnoresBlankLines(t *testing.T) {
	namesFile := bytes.NewBufferString("\nGlobal/Vim\nPython\n")
	names := ParseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf(errorTemplate, names, expectedNames)
	}
}

func TestParseNamesFileStripsSpaces(t *testing.T) {
	namesFile := bytes.NewBufferString("Global/Vim   \n  \n   Python\n")
	names := ParseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf(errorTemplate, names, expectedNames)
	}
}

func TestNamedIgnoreContentsDisplayName(t *testing.T) {
	nics := []NamedIgnoreContents{
		NamedIgnoreContents{"Vim", "*.swp"},
		NamedIgnoreContents{"Global/Vim", "*.swp"},
		NamedIgnoreContents{"Vim.gitignore", "*.swp"},
		NamedIgnoreContents{"Vim.patterns", "*.swp"},
		NamedIgnoreContents{"Global/Vim.gitignore", "*.swp"},
	}
	expectedDisplayName := "Vim"
	for _, nic := range nics {
		displayName := nic.DisplayName()
		if displayName != expectedDisplayName {
			t.Errorf(errorTemplate, displayName, expectedDisplayName)
		}
	}
}

func TestNamesToUrls(t *testing.T) {
	getter := HTTPIgnoreGetter{
		"https://raw.githubusercontent.com/github/gitignore/master",
		".gitignore",
	}
	names := []string{"Go", "Python", "Global/Vim.gitignore", "Some.patterns"}
	urls := getter.NamesToUrls(names)
	expectedURLs := []NamedSource{
		NamedSource{"Go", "https://raw.githubusercontent.com/github/gitignore/master/Go.gitignore"},
		NamedSource{"Python", "https://raw.githubusercontent.com/github/gitignore/master/Python.gitignore"},
		NamedSource{"Global/Vim.gitignore", "https://raw.githubusercontent.com/github/gitignore/master/Global/Vim.gitignore"},
		NamedSource{"Some.patterns", "https://raw.githubusercontent.com/github/gitignore/master/Some.patterns"},
	}
	if !reflect.DeepEqual(urls, expectedURLs) {
		t.Errorf(errorTemplate, urls, expectedURLs)
	}
}

func TestFailedURLsError(t *testing.T) {
	failedURLs := new(FailedSources)
	failedURLs.Add(
		&FailedSource{
			"https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore",
			fmt.Errorf("status code 404")})
	failedURLs.Add(
		&FailedSource{
			"https://raw.githubusercontent.com/github/gitignore/master/Totally.gitignore",
			fmt.Errorf("Error reading response body: too many 💩s")})
	expectedErrorStr := `Errors for the following URLs:
https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore status code 404
https://raw.githubusercontent.com/github/gitignore/master/Totally.gitignore Error reading response body: too many 💩s`
	errorStr := failedURLs.Error()
	if errorStr != expectedErrorStr {
		t.Errorf(errorTemplate, errorStr, expectedErrorStr)
	}
}

type NameWithExpectedContents struct {
	name             string
	expectedURLPath  string
	expectedContents string
	expectedError    error
}

func TestGetIgnoreFilesForNameOnly(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(handlerFunc))
	defer testServer.Close()
	assertGetIgnoreFilesReturnsExpectedContents(
		t,
		testServer,
		[]NameWithExpectedContents{
			{
				"Global/Vim",
				"/Global/Vim.gitignore",
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
		[]NameWithExpectedContents{
			{
				"Global/Vim.gitignore",
				"/Global/Vim.gitignore",
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
		[]NameWithExpectedContents{
			{
				"Foo.bar",
				"/Foo.bar",
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
		[]NameWithExpectedContents{
			{
				"Nonexistent",
				"/Nonexistent.gitignore",
				"",
				fmt.Errorf("Got status code 404"),
			},
		},
	)
}

func assertGetIgnoreFilesReturnsExpectedContents(t *testing.T, testServer *httptest.Server, namesWithContents []NameWithExpectedContents) {
	getter := HTTPIgnoreGetter{
		testServer.URL,
		".gitignore",
	}
	contentsChannel := make(chan RetrievedContents)
	var requestsPending sync.WaitGroup
	for _, nameWithContents := range namesWithContents {
		expectedContents := RetrievedContents{
			NamedSource{nameWithContents.name, testServer.URL + nameWithContents.expectedURLPath},
			nameWithContents.expectedContents,
			nameWithContents.expectedError,
		}
		go getter.GetIgnoreFiles([]string{nameWithContents.name}, contentsChannel, &requestsPending)
		gotContents := <-contentsChannel
		if !reflect.DeepEqual(gotContents, expectedContents) {
			t.Errorf(errorTemplate, gotContents, expectedContents)
		}
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

func TestWriteIgnoreFile(t *testing.T) {
	ignoreFile := bytes.NewBufferString("")
	responseContents := []NamedIgnoreContents{
		NamedIgnoreContents{name: "Global/Vim", contents: ".*.swp\ntags\n"},
		NamedIgnoreContents{name: "Go.gitignore", contents: "*.o\n*.exe\n"},
	}
	writeIgnoreFile(ignoreFile, responseContents)
	ignoreFileContents := ignoreFile.String()
	expectedContents := `#######
# Vim #
#######
.*.swp
tags


######
# Go #
######
*.o
*.exe
`
	if ignoreFileContents != expectedContents {
		t.Errorf(errorTemplate, ignoreFileContents, expectedContents)
	}
}
