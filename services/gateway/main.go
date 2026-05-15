package main

import (
	"goapp/framework/server"
	"goapp/framework/server/http"
	"log"
	nethttp "net/http"

	"github.com/labstack/echo/v5"
)

func main() {
	srv, err := http.New()

	if err != nil {
		log.Fatal(err)
	}

	srv.Echo.GET("/", func(c *echo.Context) error {
		return c.String(nethttp.StatusOK, "Hello, World!")
	})

	server.Run(srv)

	select {}
}
