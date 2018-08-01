package main

import (
	"net/http"

	"github.com/labstack/echo"
)

func main() {
	server := echo.New()
	server.GET("/", helloWorld)
	server.Start(":1234")
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello CLΛRK!")
}
