package errors

import (
	"fmt"
	"testing"

	"github.com/gotgenes/getignore/testutils"
)

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
		testutils.TError(t, errorStr, expectedErrorStr)
	}
}
