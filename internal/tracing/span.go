package tracing

import (
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// SetSpanError は err があれば Span を Error 状態にし、エラーを記録する。
// RecordError だけでは Jaeger などのバックエンドで Span が失敗と見なされないため、あわせて SetStatus も呼ぶ。
func SetSpanError(span trace.Span, err error) {
	if err == nil {
		return
	}
	span.SetStatus(codes.Error, err.Error())
	span.RecordError(err)
}
