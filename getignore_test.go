package main

import (
	"strings"
	"testing"

	"github.com/gotgenes/getignore/contentstructs"
	"github.com/gotgenes/getignore/testutils"
)

func TestParseNamesFile(t *testing.T) {
	namesFile := strings.NewReader("Global/Vim\nPython\n")
	names := ParseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	testutils.AssertDeepEqual(t, names, expectedNames)
}

func TestParseNamesFileIgnoresBlankLines(t *testing.T) {
	namesFile := strings.NewReader("\nGlobal/Vim\nPython\n")
	names := ParseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	testutils.AssertDeepEqual(t, names, expectedNames)
}

func TestParseNamesFileStripsSpaces(t *testing.T) {
	namesFile := strings.NewReader("Global/Vim   \n  \n   Python\n")
	names := ParseNamesFile(namesFile)
	expectedNames := []string{"Global/Vim", "Python"}
	testutils.AssertDeepEqual(t, names, expectedNames)
}

func TestNamedIgnoreContentsDisplayName(t *testing.T) {
	nics := []contentstructs.NamedIgnoreContents{
		{"Vim", "*.swp"},
		{"Global/Vim", "*.swp"},
		{"Vim.gitignore", "*.swp"},
		{"Vim.patterns", "*.swp"},
		{"Global/Vim.gitignore", "*.swp"},
	}
	expectedDisplayName := "Vim"
	for _, nic := range nics {
		displayName := nic.DisplayName()
		if displayName != expectedDisplayName {
			testutils.TError(t, displayName, expectedDisplayName)
		}
	}
}
