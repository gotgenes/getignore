package writers_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWriters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Writers Suite")
}
