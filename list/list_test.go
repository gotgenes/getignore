package list

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const contentsBodyTemplate = `{
  "sha": "be962308a6efb6d27446d13ee05b74948e616a01",
  "url": "https://api.github.com/repos/github/gitignore/git/trees/be962308a6efb6d27446d13ee05b74948e616a01",
  "tree": [
%v
  ],
  "truncated": false
 }`

func TestParseGitTreeToFileNamesNoTreeContents(t *testing.T) {
	assertReturnsExpectedFileNames(t, "", nil)
}

func TestParseGitTreeToFileNames(t *testing.T) {
	treeContents := `
    {
      "path": "Global/Vim.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "42e7afc100512b21af0a3339fe851d21e3ee4ed5",
      "size": 147,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/42e7afc100512b21af0a3339fe851d21e3ee4ed5"
    },
    {
      "path": "Go.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "a1338d68517ee2ad6ee11214b201e5958cb2bbc3",
      "size": 275,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/a1338d68517ee2ad6ee11214b201e5958cb2bbc3"
    }`
	assertReturnsExpectedFileNames(t, treeContents, []string{"Global/Vim.gitignore", "Go.gitignore"})
}

func TestParseGitTreeToFileNamesIgnoresDirectories(t *testing.T) {
	treeContents := `
    {
      "path": ".github",
      "mode": "040000",
      "type": "tree",
      "sha": "0b994405766547cfaeaf5e45a176947336926c23",
      "url": "https://api.github.com/repos/github/gitignore/git/trees/0b994405766547cfaeaf5e45a176947336926c23"
    },
    {
      "path": "Global/Vim.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "42e7afc100512b21af0a3339fe851d21e3ee4ed5",
      "size": 147,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/42e7afc100512b21af0a3339fe851d21e3ee4ed5"
    },
    {
      "path": "Go.gitignore",
      "mode": "100644",
      "type": "blob",
      "sha": "a1338d68517ee2ad6ee11214b201e5958cb2bbc3",
      "size": 275,
      "url": "https://api.github.com/repos/github/gitignore/git/blobs/a1338d68517ee2ad6ee11214b201e5958cb2bbc3"
    }`
	assertReturnsExpectedFileNames(t, treeContents, []string{"Global/Vim.gitignore", "Go.gitignore"})
}

func TestParseGitTreeToFileNamesReturnsDecodeError(t *testing.T) {
	assertReturnsError(
		t,
		`}"path": "Go.gitignore", "type": "blob"`,
		"invalid character '}' looking for beginning of value")
}

func TestParseGitTreeToFileNamesReturnsUnmarshallError(t *testing.T) {
	assertReturnsError(
		t,
		`{"path": "Go.gitignore", "type": 1}`,
		"json: cannot unmarshal number into Go struct field fileInfo.Type of type string")
}

func assertReturnsExpectedFileNames(t *testing.T, treeContents string, expectedFileNames []string) {
	contents := fmt.Sprintf(contentsBodyTemplate, treeContents)
	responseBody := strings.NewReader(contents)
	fileNames, _ := parseGitTreeToFileNames(responseBody)
	assert.Equal(t, expectedFileNames, fileNames)
}

func assertReturnsError(t *testing.T, treeContents string, expectedErrorMessage string) {
	contents := fmt.Sprintf(contentsBodyTemplate, treeContents)
	responseBody := strings.NewReader(contents)
	_, err := parseGitTreeToFileNames(responseBody)
	assert.EqualError(t, err, expectedErrorMessage)
}

func TestFilterBySuffix(t *testing.T) {
	assertReturnsFilteredFileNames(t, nil, "", nil)
	assertReturnsFilteredFileNames(
		t,
		[]string{"Foo", "Global/Vim.gitignore", "Go.gitignore", "Ignoreme"},
		"",
		[]string{"Foo", "Global/Vim.gitignore", "Go.gitignore", "Ignoreme"})
	assertReturnsFilteredFileNames(
		t,
		[]string{"Foo", "Global/Vim.gitignore", "Go.gitignore", "Ignoreme"},
		".gitignore",
		[]string{"Global/Vim.gitignore", "Go.gitignore"})
}

func assertReturnsFilteredFileNames(t *testing.T, fileNames []string, suffix string, expectedFileNames []string) {
	assert.Equal(t, expectedFileNames, filterBySuffix(fileNames, suffix))
}
