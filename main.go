package main

import (
	"flag"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type FileType string

const (
	Dir     = FileType("directory")
	Jpeg    = FileType("jpeg")
	Png     = FileType("png")
	Txt     = FileType("txt")
	Json    = FileType("json")
	UnKnown = FileType("unknown")
)

type Entity struct {
	Path     string
	FileType FileType
}

func newEntity(path string) *Entity {
	return &Entity{Path: path, FileType: getFileType(path)}
}

var root *string

func main() {
	// コマンドライン引数
	root = flag.String("root", "path/to/hoge", "Specify root directory path to open.")
	port := flag.Int("port", 8000, "Specify port to use.")
	flag.Parse()
	fmt.Printf("root: %s\n", *root)
	fmt.Printf("port: %d\n", *port)

	// Echo instance
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
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

// 拡張子からファイル種類を判別
func getFileType(path string) FileType {
	ext := filepath.Ext(path) // "path/to/hoge.c" => ".c"
	fmt.Printf("ext: %s\n", ext)

	switch ext {
	case "":
		return Dir
	case ".jpeg":
		return Jpeg
	case ".jpg":
		return Jpeg
	case ".png":
		return Png
	case ".txt":
		return Txt
	case ".json":
		return Json
	default:
		return UnKnown
	}
}
