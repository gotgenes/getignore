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
		{Name: "Global/Vim", Contents: ".*.swp\ntags\n"},
		{Name: "Go.gitignore", Contents: "*.o\n*.exe\n"},
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

func TestWriteIgnoreFileEnsuresTerminatingNewlines(t *testing.T) {
	ignoreFile := bytes.NewBufferString("")
	responseContents := []contentstructs.NamedIgnoreContents{
		{Name: "Global/Vim", Contents: ".*.swp\ntags"},
		{Name: "Go.gitignore", Contents: "*.o\n*.exe"},
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
