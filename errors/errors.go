package errors

import (
	"fmt"
	"strings"
)

// FailedSource represents a source unable to be retrieved or processed
type FailedSource struct {
	Source string
	Err    error
}

func (fs *FailedSource) Error() string {
	return fmt.Sprintf("%s: %s", fs.Source, fs.Err.Error())
}

// FailedSources represents a collection of FailedSource instances
type FailedSources []FailedSource

func (failedSources FailedSources) Error() string {
	sourceErrors := make([]string, len(failedSources))
	for i, failedSource := range failedSources {
		sourceErrors[i] = failedSource.Error()
	}
	stringOfErrors := strings.Join(sourceErrors, "\n")
	return "Errors retrieving the following sources:\n" + stringOfErrors
}
