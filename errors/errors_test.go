package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFailedSourcesError(t *testing.T) {
	failedSources := FailedSources{
		FailedSource{
			"https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore",
			fmt.Errorf("status code 404"),
		},
		FailedSource{
			"https://raw.githubusercontent.com/github/gitignore/master/Totally.gitignore",
			fmt.Errorf("Error reading response body: too many 💩s"),
		},
	}
	expectedErrorStr := `Errors retrieving the following sources:
https://raw.githubusercontent.com/github/gitignore/master/Bogus.gitignore: status code 404
https://raw.githubusercontent.com/github/gitignore/master/Totally.gitignore: Error reading response body: too many 💩s`
	errorStr := failedSources.Error()
	assert.Equal(t, expectedErrorStr, errorStr)
}
