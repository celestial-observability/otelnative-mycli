package tracing

import (
	"context"
	"os"
)

// Init はトレース出力を初期化し、終了処理用の関数を返す。
// エクスポーターは環境変数 HI_TRACES_EXPORTER で選ぶ。この節では none（デフォルト）のみ対応する。
func Init(_ context.Context, _ string) (func(context.Context) error, error) {
	exporterName := os.Getenv("HI_TRACES_EXPORTER")
	if exporterName == "" || exporterName == "none" {
		return func(context.Context) error { return nil }, nil
	}

	// エクスポーターの生成と TracerProvider の組み立ては次節で実装する。
	return func(context.Context) error { return nil }, nil
}
