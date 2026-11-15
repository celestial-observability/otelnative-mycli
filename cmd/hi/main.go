// hi は CLI ツールのエントリーポイント
package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/celestial-observability/otelnative-mycli/internal/cmd"
)

var CLI struct {
	Ja cmd.JaCmd `cmd:"ja" help:"日本語で挨拶を出力する"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Hello, World!")
		return
	}
	kongCtx := kong.Parse(&CLI)
	if err := kongCtx.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
