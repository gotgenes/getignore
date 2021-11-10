package errors_test

import (
	"errors"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	gierrors "github.com/gotgenes/getignore/errors"
)

var _ = Describe("FailedFile", func() {
	var ff gierrors.FailedFile

	BeforeEach(func() {
		ff = gierrors.FailedFile{
			Name:    "Go.gitignore",
			Message: "could not connect",
			Err:     fmt.Errorf("problem connecting to the server"),
		}
	})

	It("should implement the error interface", func() {
		Expect(ff).Should(MatchError("failed to get Go.gitignore: could not connect"))
	})

	It("should support unwrapping the inner error", func() {
		Expect(errors.Unwrap(ff)).Should(MatchError("problem connecting to the server"))
	})
})
