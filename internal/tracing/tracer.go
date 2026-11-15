// Package tracing は hi コマンドのトレース出力を担う。
package tracing

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Tracer は計装で使う共通のトレーサーを返す。名前を一箇所で管理するため、各呼び出し元はこの関数を経由する。
func Tracer() trace.Tracer {
	return otel.Tracer("hi")
}
