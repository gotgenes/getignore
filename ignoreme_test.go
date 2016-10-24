package ignoreme

import "testing"

var error_template string = "Got %q, expected %q"

func TestIgnoreFetcher(t *testing.T) {
	baseUrl := "https://github.com/github/gitignore"
	fetcher := ignoreFetcher{baseUrl: baseUrl}
	gotUrl := fetcher.baseUrl
	if gotUrl != baseUrl {
		t.Errorf(error_template, gotUrl, baseUrl)
	}
}

func TestNameToUrl(t *testing.T) {
	fetcher := ignoreFetcher{baseUrl: "https://github.com/github/gitignore"}
	url := fetcher.NameToUrl("Go")
	expectedUrl := "https://github.com/github/gitignore/Go.gitignore"
	if url != expectedUrl {
		t.Errorf(error_template, url, expectedUrl)
	}
}
