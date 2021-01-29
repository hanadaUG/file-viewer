package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// コマンドライン引数
	root := flag.String("root", "path/to/hoge", "Specify root directory path to open.")
	port := flag.Int("port", 8000, "Specify port to use.")
	flag.Parse()
	fmt.Printf("root: %s\n", *root)
	fmt.Printf("port: %d\n", *port)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", handler)

	// Start server
	address := fmt.Sprintf(":%d", *port)
	e.Logger.Fatal(e.Start(address))
}

// Handler
func handler(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
