package writers

import (
	"bytes"
	"testing"

	"github.com/gotgenes/getignore/contentstructs"
	"github.com/stretchr/testify/assert"
)

func TestWriteIgnoreFile(t *testing.T) {
	retrievedContents := []contentstructs.NamedIgnoreContents{
		{Name: "Global/Vim", Contents: ".*.swp\ntags\n"},
		{Name: "Go.gitignore", Contents: "*.o\n*.exe\n"},
	}
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
	assertWritesExpectedContents(t, retrievedContents, expectedContents)
}

func TestWriteIgnoreFileEnsuresTerminatingNewlines(t *testing.T) {
	retrievedContents := []contentstructs.NamedIgnoreContents{
		{Name: "Global/Vim", Contents: ".*.swp\ntags"},
		{Name: "Go.gitignore", Contents: "*.o\n*.exe"},
	}
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
	assertWritesExpectedContents(t, retrievedContents, expectedContents)
}

func TestWriteIgnoreFileNoContents(t *testing.T) {
	retrievedContents := []contentstructs.NamedIgnoreContents{
		{Name: "Global/Vim", Contents: ""},
		{Name: "Go.gitignore", Contents: "\n"},
	}
	expectedContents := `#######
# Vim #
#######


######
# Go #
######
`
	assertWritesExpectedContents(t, retrievedContents, expectedContents)
}

func TestWriteIgnoreFileStripsLeadingAndTrailingWhitespace(t *testing.T) {
	retrievedContents := []contentstructs.NamedIgnoreContents{
		{Name: "Global/Vim", Contents: "\n    \n.*.swp\ntags\n"},
		{Name: "Go.gitignore", Contents: "*.o\n*.exe     \n\n\t\n"},
	}
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
	assertWritesExpectedContents(t, retrievedContents, expectedContents)
}

func assertWritesExpectedContents(t *testing.T, retrievedContents []contentstructs.NamedIgnoreContents, expectedContents string) {
	ignoreFile := bytes.NewBufferString("")
	WriteIgnoreFile(ignoreFile, retrievedContents)
	ignoreFileContents := ignoreFile.String()
	assert.Equal(t, expectedContents, ignoreFileContents)
}
