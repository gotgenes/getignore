package contentstructs

import (
	"path/filepath"
	"strings"
)

// NamedContents represents the contents (patterns and comments) of a
// gitignore file
type NamedContents struct {
	Name     string
	Contents string
}

// DisplayName returns the decorated name, suitable for a section header in a
// gitignore file
func (n *NamedContents) DisplayName() string {
	baseName := filepath.Base(n.Name)
	return strings.TrimSuffix(baseName, filepath.Ext(baseName))
}
