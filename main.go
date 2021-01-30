package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
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

type Template struct {
	templates *template.Template
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Register a static file
	e.Static("/assets", "assets")

	// Template
	funcMap := template.FuncMap{
		//"add":  func(a, b int) int { return a + b },
	}

	tmpl, err := template.New("t").Funcs(funcMap).ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	t := &Template{
		templates: tmpl,
	}

	e.Renderer = t

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

	//埋め込み変数
	data := map[string]interface{}{
		"entity":   entity,
		"entities": entities,
	}

	return c.Render(http.StatusOK, "index.html", data)
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
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

// template から呼び出すため public なメソッドにする
func (entity Entity) IsDir() bool {
	return entity.FileType == Dir
}

func (entity Entity) GetParentDir() string {
	return path.Dir(entity.Path)
}
