package contentstructs

import (
	"testing"

	"github.com/gotgenes/getignore/testutils"
)

func TestNamedIgnoreContentsDisplayName(t *testing.T) {
	nics := []NamedIgnoreContents{
		{Name: "Vim", Contents: "*.swp"},
		{Name: "Global/Vim", Contents: "*.swp"},
		{Name: "Vim.gitignore", Contents: "*.swp"},
		{Name: "Vim.patterns", Contents: "*.swp"},
		{Name: "Global/Vim.gitignore", Contents: "*.swp"},
	}
	expectedDisplayName := "Vim"
	for _, nic := range nics {
		displayName := nic.DisplayName()
		if displayName != expectedDisplayName {
			testutils.TError(t, displayName, expectedDisplayName)
		}
	}
}
