package getignore

import (
	"bufio"
	"io"
	"strings"
)

// ParseNamesFile reads a file containing one name of a gitignore patterns file per line
func ParseNamesFile(namesFile io.Reader) []string {
	var a []string
	scanner := bufio.NewScanner(namesFile)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if len(name) > 0 {
			a = append(a, name)
		}
	}
	return a
}
