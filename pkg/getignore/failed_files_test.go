package getignore_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/gotgenes/getignore/pkg/getignore"
)

var _ = Describe("FailedFiles", func() {
	It("should implement the error interface", func() {
		err := getignore.FailedFiles{
			{
				Name:    "Go.gitignore",
				Message: "could not connect",
				Err:     fmt.Errorf("problem connecting to the server"),
			},
			{
				Name:    "Nonexistent.gitignore",
				Message: "file not found in tree",
			},
		}
		expectedMsg := `failed to get the following files: Go.gitignore, Nonexistent.gitignore
Go.gitignore: could not connect
Nonexistent.gitignore: file not found in tree
`
		Expect(err).Should(MatchError(expectedMsg))
	})
})
