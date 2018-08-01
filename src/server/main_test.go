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
	Context(".fancyAdd", func() {
		It("Returns the length of the given parameter plus 42", func() {
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(httptest.NewRequest(echo.GET, "/", strings.NewReader("")), rec)
			c.SetParamNames("value")
			c.SetParamValues("ClarkKent")
			Expect(fancyAdd(c)).To(BeNil())
			Expect(rec.Code).To(Equal(http.StatusOK))
			Expect(rec.Body.String()).To(Equal("The result is 51"))
		})
	})
})
