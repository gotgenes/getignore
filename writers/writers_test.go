package writers

import (
	"bytes"
	"testing"

	"github.com/gotgenes/getignore/contentstructs"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, expectedContents, ignoreFileContents)
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
		assert.Equal(t, expectedContents, ignoreFileContents)
	}
}
