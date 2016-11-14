package main

import (
	"bytes"
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

func TestFailedURLsAdd(t *testing.T) {
	failedURLs := new(FailedURLs)
	failedURLs.Add("https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore")
	failedURLs.Add("https://raw.githubusercontent.com/github/gitignore/master/Totally.gitignore")
	expectedURLs := []string{
		"https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore",
		"https://raw.githubusercontent.com/github/gitignore/master/Totally.gitignore",
	}
	if !reflect.DeepEqual(failedURLs.URLs, expectedURLs) {
		t.Errorf(errorTemplate, failedURLs.URLs, expectedURLs)
	}
}

func TestFailedURLsError(t *testing.T) {
	failedURLs := new(FailedURLs)
	failedURLs.Add("https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore")
	failedURLs.Add("https://raw.githubusercontent.com/github/gitignore/master/Totally.gitignore")
	expectedErrorStr := `Failed to retrieve or read content from the following URLs:
https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore
https://raw.githubusercontent.com/github/gitignore/master/Totally.gitignore`
	errorStr := failedURLs.Error()
	if errorStr != expectedErrorStr {
		t.Errorf(errorTemplate, errorStr, expectedErrorStr)
	}
}

func TestWriteIgnoreFile(t *testing.T) {
	ignoreFile := bytes.NewBufferString("")
	responseContents := []NamedIgnoreContents{
		NamedIgnoreContents{name: "Vim", contents: ".*.swp\ntags\n"},
		NamedIgnoreContents{name: "Go", contents: "*.o\n*.exe\n"},
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
