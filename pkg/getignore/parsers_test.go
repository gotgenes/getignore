package getignore_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gotgenes/getignore/pkg/getignore"
)

var _ = Describe("ParseNamesFile", func() {
	assertReturnsExpectedNames := func(contents string) {
		namesFile := strings.NewReader(contents)
		names := getignore.ParseNamesFile(namesFile)
		Expect(names).Should(Equal([]string{"Global/Vim", "Python"}))
	}

	It("parses a standard file", func() {
		assertReturnsExpectedNames("Global/Vim\nPython\n")
	})

	It("ignores blank lines", func() {
		assertReturnsExpectedNames("\nGlobal/Vim\n\nPython\n\n")
	})

	It("strips whitespace", func() {
		assertReturnsExpectedNames("Global/Vim   \n  \n   Python\n")
	})
})
