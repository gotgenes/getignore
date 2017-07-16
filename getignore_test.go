package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNamesFile(t *testing.T) {
	assertReturnsExpectedNames(t, "Global/Vim\nPython\n", []string{"Global/Vim", "Python"})
}

func TestParseNamesFileIgnoresBlankLines(t *testing.T) {
	assertReturnsExpectedNames(t, "\nGlobal/Vim\nPython\n", []string{"Global/Vim", "Python"})
}

func TestParseNamesFileStripsSpaces(t *testing.T) {
	assertReturnsExpectedNames(
		t, "Global/Vim   \n  \n   Python\n", []string{"Global/Vim", "Python"})
}

func assertReturnsExpectedNames(t *testing.T, namesFileContents string, expectedNames []string) {
	namesFile := strings.NewReader(namesFileContents)
	names := ParseNamesFile(namesFile)
	assert.Equal(t, expectedNames, names)
}
