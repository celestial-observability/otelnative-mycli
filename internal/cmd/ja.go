// Package cmd はCLIサブコマンド層を実装する
package cmd

import (
	"context"
	"fmt"

	"github.com/celestial-observability/otelnative-mycli/internal/tracing"
)

// JaCmd は ja サブコマンドを実装する
type JaCmd struct{}

// Run は標準出力に「こんにちは、世界!」を出力する
func (c *JaCmd) Run(ctx context.Context) error {
	_, span := tracing.Tracer().Start(ctx, "ja")
	defer span.End()
	fmt.Println("こんにちは、世界!")
	return nil
}
