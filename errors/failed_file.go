package errors

import (
	"fmt"
	"strings"
)

// FailedFile represents a gitignore file unable to be retrieved or processed
type FailedFile struct {
	Name    string
	Message string
	Err     error
}

func (f FailedFile) Error() string {
	return fmt.Sprintf("failed to get %s: %s", f.Name, f.Message)
}

func (f FailedFile) Unwrap() error {
	return f.Err
}

// FailedFiles represents a collection of FailedFile instances
type FailedFiles []FailedFile

func (e FailedFiles) Error() string {
	fileNames := make([]string, len(e))
	reasons := make([]string, len(e))
	for i, failedFile := range e {
		fileNames[i] = failedFile.Name
		reasons[i] = fmt.Sprintf("%s: %s", failedFile.Name, failedFile.Message)
	}
	filesStr := strings.Join(fileNames, ", ")
	reasonsStr := strings.Join(reasons, "\n")
	return fmt.Sprintf("failed to get the following files: %s\n%s\n", filesStr, reasonsStr)
}
