package contentstructs

import (
	"path/filepath"
	"strings"
)

// NamedIgnoreContents represents the contents (patterns and comments) of a
// gitignore file
type NamedIgnoreContents struct {
	Name     string
	Contents string
}

// DisplayName returns the decorated name, suitable for a section header in a
// gitignore file
func (nic *NamedIgnoreContents) DisplayName() string {
	baseName := filepath.Base(nic.Name)
	return strings.TrimSuffix(baseName, filepath.Ext(baseName))
}
