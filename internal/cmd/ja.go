// Package cmd はCLIサブコマンド層を実装する
package cmd

import "fmt"

// JaCmd は ja サブコマンドを実装する
type JaCmd struct{}

// Run は標準出力に「こんにちは、世界!」を出力する
func (c *JaCmd) Run() error {
	fmt.Println("こんにちは、世界!")
	return nil
}
