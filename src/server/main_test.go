package main

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/labstack/echo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Main", func() {
	Context(".helloWorld", func() {
		It("Returns a String containing 'Hello World'", func() {
			e := echo.New()
			req := httptest.NewRequest(echo.GET, "/", strings.NewReader(""))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			Expect(helloWorld(c)).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(Equal("Hello CLΛRK!"))
		})
	})
})
