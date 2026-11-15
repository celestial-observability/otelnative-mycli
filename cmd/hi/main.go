// hi は CLI ツールのエントリーポイント
package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/alecthomas/kong"

	"github.com/celestial-observability/otelnative-mycli/internal/cmd"
	"github.com/celestial-observability/otelnative-mycli/internal/tracing"
)

var CLI struct {
	Ja cmd.JaCmd `cmd:"ja" help:"日本語で挨拶を出力する"`
}

func main() {
	os.Exit(run())
}

func run() int {
	ctx := context.Background()
	shutdown, err := tracing.Init(ctx, "hi")
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	defer func() { _ = shutdown(ctx) }()

	ctx, span := tracing.Tracer().Start(ctx, "main")
	defer span.End()

	if len(os.Args) == 1 {
		fmt.Println("Hello, World!")
		return 0
	}
	kongCtx := kong.Parse(&CLI,
		kong.Exit(func(code int) {
			if code != 0 {
				tracing.SetSpanError(span, errors.New("kong: parse error"))
			}
			span.End()
			_ = shutdown(ctx)
			os.Exit(code)
		}),
		kong.BindFor[context.Context](ctx),
	)
	if err := kongCtx.Run(); err != nil {
		tracing.SetSpanError(span, err)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
