// hi は "Hello World" を標準出力するプログラム
package main

import (
	"fmt"

	"github.com/alecthomas/kong"
)

var CLI struct{}

func main() {
	kong.Parse(&CLI)
	fmt.Println("Hello, World!")
}
