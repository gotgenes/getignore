package main

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

var errorTemplate = "Got %q, expected %q"

func TestParseNamesFile(t *testing.T) {
	namesFile := bytes.NewBufferString("Global/Vim\nPython\n")
	names := parseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf(errorTemplate, names, expectedNames)
	}
}

func TestParseNamesFileIgnoresBlankLines(t *testing.T) {
	namesFile := bytes.NewBufferString("\nGlobal/Vim\nPython\n")
	names := parseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf(errorTemplate, names, expectedNames)
	}
}

func TestParseNamesFileStripsSpaces(t *testing.T) {
	namesFile := bytes.NewBufferString("Global/Vim   \n  \n   Python\n")
	names := parseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf(errorTemplate, names, expectedNames)
	}
}

func TestIgnoreFetcher(t *testing.T) {
	baseURL := "https://github.com/github/gitignore"
	fetcher := IgnoreFetcher{baseURL: baseURL}
	gotURL := fetcher.baseURL
	if gotURL != baseURL {
		t.Errorf(errorTemplate, gotURL, baseURL)
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
	fetcher := IgnoreFetcher{baseURL: "https://raw.githubusercontent.com/github/gitignore/master"}
	names := []string{"Go", "Python"}
	urls := fetcher.NamesToUrls(names)
	expectedURLs := []NamedURL{
		NamedURL{"Go", "https://raw.githubusercontent.com/github/gitignore/master/Go.gitignore"},
		NamedURL{"Python", "https://raw.githubusercontent.com/github/gitignore/master/Python.gitignore"},
	}
	if !reflect.DeepEqual(urls, expectedURLs) {
		t.Errorf(errorTemplate, urls, expectedURLs)
	}
}

func TestNameToUrl(t *testing.T) {
	fetcher := IgnoreFetcher{baseURL: "https://github.com/github/gitignore"}
	url := fetcher.NameToURL("Go")
	expectedURL := NamedURL{"Go", "https://github.com/github/gitignore/Go.gitignore"}
	if url != expectedURL {
		t.Errorf(errorTemplate, url, expectedURL)
	}
}

func TestFailedURLsError(t *testing.T) {
	failedURLs := new(FailedURLs)
	failedURLs.Add(
		&FailedURL{
			"https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore",
			fmt.Errorf("status code 404")})
	failedURLs.Add(
		&FailedURL{
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
