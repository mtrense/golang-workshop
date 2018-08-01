package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	server := echo.New()
	server.GET("/", helloWorld)
	server.GET("/fancyadd/:value", fancyAdd)
	server.Start(":1234")
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello CLΛRK!")
}

func fancyAdd(c echo.Context) error {
	value := c.Param("value")
	return c.String(http.StatusOK, fmt.Sprintf("The result is %d\n", len(value)+42))
}
