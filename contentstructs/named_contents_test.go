package contentstructs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/gotgenes/getignore/contentstructs"
)

var _ = Describe("NamedContents", func() {
	Describe("DisplayName", func() {
		It("should return the name", func() {
			nc := contentstructs.NamedContents{Name: "Vim"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})

		It("should take only the base name", func() {
			nc := contentstructs.NamedContents{Name: "Global/Vim"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})

		It("should strip the extension", func() {
			nc := contentstructs.NamedContents{Name: "Vim.gitignore"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})

		It("should strip any extension", func() {
			nc := contentstructs.NamedContents{Name: "Vim.extension"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})

		It("should take only the base name and strip the extension", func() {
			nc := contentstructs.NamedContents{Name: "Global/Vim.gitignore"}
			Expect(nc.DisplayName()).Should(Equal("Vim"))
		})
	})
})
