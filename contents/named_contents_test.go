package contents_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gotgenes/getignore/contents"
)

var _ = Describe("NamedContents", func() {
	Describe("DisplayName", func() {
		It("should return the name", func() {
			nc := contents.NamedContents{Name: "Vim"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})

		It("should take only the base name", func() {
			nc := contents.NamedContents{Name: "Global/Vim"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})

		It("should strip the extension", func() {
			nc := contents.NamedContents{Name: "Vim.gitignore"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})

		It("should strip any extension", func() {
			nc := contents.NamedContents{Name: "Vim.extension"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})

		It("should take only the base name and strip the extension", func() {
			nc := contents.NamedContents{Name: "Global/Vim.gitignore"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})
	})
})
