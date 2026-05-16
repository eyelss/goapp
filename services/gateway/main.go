package main

import (
	"goapp/framework/server"
	"goapp/framework/server/http"
	"log"
	http2 "net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v5"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func handleRoot(c *echo.Context) error {
	log.Println("Request root")

	b, err := os.ReadFile("./dist/index.html")

	if err != nil {
		log.Printf("Error reading file: %v", err)
	}

	return c.Blob(http2.StatusOK, "text/html", b)
}

func handleResource(c *echo.Context) error {
	path := c.Request().URL.Path
	log.Printf("Request resource %s\n", path)
	relativePath := "./dist" + path

	ext := filepath.Ext(path)

	var contentType string
	switch ext {
	case ".css":
		contentType = "text/css"
	case ".html":
		contentType = "text/html"
	case ".js":
		contentType = "application/javascript"
	case ".svg":
		contentType = "image/svg+xml"
	case ".png":
		contentType = "image/png"
	case ".jpg":
	case ".jpeg":
		contentType = "image/jpg"
	}

	if contentType == "" {
		e := &ErrorResponse{
			Message: "content type is not valid",
		}

		return c.JSON(http2.StatusNotAcceptable, e)
	}

	f, err := os.Open(relativePath)

	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			e := &ErrorResponse{
				Message: "resource not found",
			}

			return c.JSON(http2.StatusNotFound, e)
		}
		log.Printf("Error reading file: %v", err)
		return err
	}

	defer f.Close()

	return c.Stream(http2.StatusOK, contentType, f)
}

func main() {
	srv, err := http.New()

	if err != nil {
		log.Fatal(err)
	}

	srv.Echo.GET("/", handleRoot)

	srv.Echo.GET("/*", handleResource)

	server.Run(srv)

	select {}
}
