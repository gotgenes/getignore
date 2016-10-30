package main

import (
	"bytes"
	"reflect"
	"testing"
)

var error_template string = "Got %q, expected %q"

func TestParseNamesFile(t *testing.T) {
	namesFile := bytes.NewBufferString("Global/Vim\nPython\n")
	names := parseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	if !reflect.DeepEqual(names, expectedNames) {
		t.Errorf(error_template, names, expectedNames)
	}
}

func TestIgnoreFetcher(t *testing.T) {
	baseUrl := "https://github.com/github/gitignore"
	fetcher := ignoreFetcher{baseUrl: baseUrl}
	gotUrl := fetcher.baseUrl
	if gotUrl != baseUrl {
		t.Errorf(error_template, gotUrl, baseUrl)
	}
}

func TestNamesToUrls(t *testing.T) {
	fetcher := ignoreFetcher{baseUrl: "https://github.com/github/gitignore"}
	names := []string{"Go", "Python"}
	namesChannel := make(chan string)
	go arrayToChannel(namesChannel, names)
	urlsChannel := make(chan string)
	go fetcher.NamesToUrls(namesChannel, urlsChannel)
	urls := channelToArray(urlsChannel)
	expectedUrls := []string{
		"https://github.com/github/gitignore/Go.gitignore",
		"https://github.com/github/gitignore/Python.gitignore",
	}
	if !reflect.DeepEqual(urls, expectedUrls) {
		t.Errorf(error_template, urls, expectedUrls)
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
	fetcher := ignoreFetcher{baseUrl: "https://github.com/github/gitignore"}
	url := fetcher.NameToUrl("Go")
	expectedUrl := "https://github.com/github/gitignore/Go.gitignore"
	if url != expectedUrl {
		t.Errorf(error_template, url, expectedUrl)
	}
}
