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
	return fmt.Sprintf("%s %s", fs.Source, fs.Err.Error())
}

// FailedSources represents a collection of FailedSource instances
type FailedSources struct {
	Sources []*FailedSource
}

// Add adds a FailedSource instance to the FailedSources collection
func (failedSources *FailedSources) Add(failedSource *FailedSource) {
	failedSources.Sources = append(failedSources.Sources, failedSource)
}

func (failedSources *FailedSources) Error() string {
	sourceErrors := make([]string, len(failedSources.Sources))
	for i, failedSource := range failedSources.Sources {
		sourceErrors[i] = failedSource.Error()
	}
	stringOfErrors := strings.Join(sourceErrors, "\n")
	return "Errors for the following URLs:\n" + stringOfErrors
}
