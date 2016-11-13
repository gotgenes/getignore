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
	fetcher := ignoreFetcher{baseURL: baseURL}
	gotURL := fetcher.baseURL
	if gotURL != baseURL {
		t.Errorf(errorTemplate, gotURL, baseURL)
	}
}

func TestNamesToUrls(t *testing.T) {
	fetcher := ignoreFetcher{baseURL: "https://raw.githubusercontent.com/github/gitignore/master"}
	names := []string{"Go", "Python"}
	urls := fetcher.NamesToUrls(names)
	expectedUrls := []string{
		"https://raw.githubusercontent.com/github/gitignore/master/Go.gitignore",
		"https://raw.githubusercontent.com/github/gitignore/master/Python.gitignore",
	}
	if !reflect.DeepEqual(urls, expectedUrls) {
		t.Errorf(errorTemplate, urls, expectedUrls)
	}
}

func arrayToChannel(c chan string, a []string) {
	for _, v := range a {
		c <- v
	}
	close(c)
}

func channelToArray(c chan string) []string {
	var a []string
	for v := range c {
		a = append(a, v)
	}
	return a
}

func TestNameToUrl(t *testing.T) {
	fetcher := ignoreFetcher{baseURL: "https://github.com/github/gitignore"}
	url := fetcher.NameToURL("Go")
	expectedURL := "https://github.com/github/gitignore/Go.gitignore"
	if url != expectedURL {
		t.Errorf(errorTemplate, url, expectedURL)
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
	expectedContents := `.*.swp
tags
*.o
*.exe
`
	if ignoreFileContents != expectedContents {
		t.Errorf(errorTemplate, ignoreFileContents, expectedContents)
	}
}
