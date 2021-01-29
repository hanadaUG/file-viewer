package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	Name     string
	Path     string
	FileType FileType
}

func newEntity(path string) *Entity {
	_, name := filepath.Split(path)
	return &Entity{Name: name, Path: path, FileType: getFileType(path)}
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
	e.GET("/*", handler)

	// Start server
	address := fmt.Sprintf(":%d", *port)
	e.Logger.Fatal(e.Start(address))
}

// Handler
func handler(c echo.Context) error {
	// root = ${HOME}/path/to/hoge
	// c.Request().URL.Path = root からの相対パス
	fullPath := filepath.Join(*root, c.Request().URL.Path)
	// fmt.Printf("fullPath: %s\n", fullPath)

	// ファイルの存在チェック
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// fmt.Printf("NotFound\n")
		return c.NoContent(http.StatusNotFound)
	} else {
		entity := newEntity(c.Request().URL.Path)
		// fmt.Printf("entity: %v\n", entity)
		return entity.Handler(c)
	}
}

func (entity Entity) Handler(c echo.Context) error {
	switch entity.FileType {
	case Dir:
		return entity.DirHandler(c)
	default:
		return entity.FileHandler(c)
	}
}

func (entity Entity) DirHandler(c echo.Context) error {
	fullPath := filepath.Join(*root, entity.Path)

	// ディレクトリ読み込み
	files, err := ioutil.ReadDir(fullPath)
	if err != nil {
		// middleware.Recover() が有効になっているため
		// サーバーは panic で終了することはない
		panic(err)
	}

	// ディレクトリ配下のファイル情報を詰め込む
	var entities []Entity
	for _, file := range files {
		path := filepath.Join(entity.Path, file.Name())
		entities = append(entities, *newEntity(path))
	}

	html := `
<html>
	<body>
		<h1>{{.entity.Path}}</h1>
{{range .entities}}
<li><a href="{{.Path}}">{{.Name}}</a></li>
{{end}}
	</body>
</html>
`
	body := bytes.NewBuffer([]byte(""))
	tmpl, err := template.New("index").Parse(html)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(body, map[string]interface{}{
		"entity":   entity,
		"entities": entities,
	})
	if err != nil {
		panic(err)
	}

	return c.HTMLBlob(http.StatusOK, body.Bytes())
}

func (entity Entity) FileHandler(c echo.Context) error {
	fullPath := filepath.Join(*root, entity.Path)

	// ファイル読み込み
	data, err := ioutil.ReadFile(fullPath)
	if err != nil {
		// middleware.Recover() が有効になっているため
		// サーバーは panic で終了することはない
		panic(err)
	}

	// ファイル種類によて MIME タイプを変更しブラウザで閲覧できるようにしている
	// 未定義のファイル種類はブラウザのダウンロード機能が走る
	var contentType string
	switch entity.FileType {
	case Jpeg:
		contentType = "image/jpeg"
	case Png:
		contentType = "image/png"
	case Txt:
		contentType = "txt/plain"
	case Json:
		contentType = "application/json"
	case UnKnown:
		fallthrough
	default:
		contentType = "application/octet-stream"
	}
	return c.Blob(http.StatusOK, contentType, data)
}

// 拡張子からファイル種類を判別
func getFileType(path string) FileType {
	ext := filepath.Ext(path) // "path/to/hoge.c" => ".c"
	ext = strings.ToLower(ext)
	// fmt.Printf("ext: %s\n", ext)

	switch ext {
	case "":
		return Dir
	case ".jpeg":
		fallthrough
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
