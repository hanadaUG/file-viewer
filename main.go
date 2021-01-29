package main

import (
	"flag"
	"fmt"
)

func main() {
	// コマンドライン引数
	{
		root := flag.String("root", "path/to/hoge", "Specify root directory path to open.")
		port := flag.Int("port", 8000, "Specify port to use.")
		flag.Parse()
		fmt.Printf("root: %s\n", *root)
		fmt.Printf("port: %d\n", *port)
	}
}
