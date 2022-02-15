package getignore_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGetignore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Getignore Suite")
}
