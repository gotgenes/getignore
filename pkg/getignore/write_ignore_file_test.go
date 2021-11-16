package getignore_test

import (
	"bytes"

	"github.com/gotgenes/getignore/pkg/getignore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WriteIgnoreFile", func() {
	var outputFile *bytes.Buffer

	BeforeEach(func() {
		outputFile = bytes.NewBufferString("")
	})

	It("should handle empty contents", func() {
		ncs := []getignore.NamedContents{
			{Name: "Global/Vim", Contents: ""},
			{Name: "Go.gitignore", Contents: "\n"},
		}
		getignore.WriteIgnoreFile(outputFile, ncs)

		expectedContents := `#######
# Vim #
#######


######
# Go #
######
`
		Expect(outputFile.String()).Should(Equal(expectedContents))
	})

	It("should write formatted contents", func() {
		ncs := []getignore.NamedContents{
			{Name: "Global/Vim", Contents: "\n    \n.*.swp\ntags\n"},
			{Name: "Go.gitignore", Contents: "*.o\n*.exe     \n\n\t\n"},
		}
		getignore.WriteIgnoreFile(outputFile, ncs)

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
		Expect(outputFile.String()).Should(Equal(expectedContents))
	})

	It("should ensure the file ends with a newline", func() {
		ncs := []getignore.NamedContents{
			{Name: "Global/Vim", Contents: ".*.swp\ntags"},
			{Name: "Go.gitignore", Contents: "*.o\n*.exe"},
		}
		getignore.WriteIgnoreFile(outputFile, ncs)

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
		Expect(outputFile.String()).Should(Equal(expectedContents))
	})

	It("should strip leading and trailing whitespace", func() {
		ncs := []getignore.NamedContents{
			{Name: "Global/Vim", Contents: "\n    \n.*.swp\ntags\n"},
			{Name: "Go.gitignore", Contents: "*.o\n*.exe     \n\n\t\n"},
		}
		getignore.WriteIgnoreFile(outputFile, ncs)

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
		Expect(outputFile.String()).Should(Equal(expectedContents))
	})
})
