package contentstructs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNamedIgnoreContentsDisplayName(t *testing.T) {
	nics := []NamedIgnoreContents{
		{Name: "Vim", Contents: "*.swp"},
		{Name: "Global/Vim", Contents: "*.swp"},
		{Name: "Vim.gitignore", Contents: "*.swp"},
		{Name: "Vim.patterns", Contents: "*.swp"},
		{Name: "Global/Vim.gitignore", Contents: "*.swp"},
	}
	for _, nic := range nics {
		displayName := nic.DisplayName()
		assert.Equal(t, "Vim", displayName)
	}
}
