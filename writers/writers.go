package writers

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/gotgenes/getignore/contentstructs"
)

// WriteIgnoreFile writes contents to a gitignore file
func WriteIgnoreFile(ignoreFile io.Writer, contents []contentstructs.NamedIgnoreContents) (err error) {
	writer := bufio.NewWriter(ignoreFile)
	for i, nc := range contents {
		if i > 0 {
			writer.WriteString("\n\n")
		}
		writer.WriteString(decorateName(nc.DisplayName()))
		writer.WriteString(nc.Contents)
		if !strings.HasSuffix(nc.Contents, "\n") {
			writer.WriteString("\n")
		}
	}
	if writer.Flush() != nil {
		err = writer.Flush()
	}
	return
}

func decorateName(name string) string {
	nameLength := len(name)
	fullHashLine := strings.Repeat("#", nameLength+4)
	nameLine := fmt.Sprintf("# %s #", name)
	decoratedName := strings.Join([]string{fullHashLine, nameLine, fullHashLine, ""}, "\n")
	return decoratedName
}
