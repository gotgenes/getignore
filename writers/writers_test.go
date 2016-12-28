package writers

import (
	"bytes"
	"testing"

	"github.com/gotgenes/getignore/contentstructs"
	"github.com/gotgenes/getignore/testutils"
)

func TestWriteIgnoreFile(t *testing.T) {
	ignoreFile := bytes.NewBufferString("")
	responseContents := []contentstructs.NamedIgnoreContents{
		{"Global/Vim", ".*.swp\ntags\n"},
		{"Go.gitignore", "*.o\n*.exe\n"},
	}
	WriteIgnoreFile(ignoreFile, responseContents)
	ignoreFileContents := ignoreFile.String()
	expectedContents := `#######
# Vim #
#######
.*.swp
tags


######
# Go #
######
*.o
*.exe
`
	if ignoreFileContents != expectedContents {
		testutils.TError(t, ignoreFileContents, expectedContents)
	}
}
