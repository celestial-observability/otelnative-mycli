package tracing

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
)

// Init はトレース出力を初期化し、終了処理用の関数を返す。
// エクスポーターは環境変数 HI_TRACES_EXPORTER（none | stdout | otlp）で選ぶ。デフォルトは none。
// otlp の宛先は環境変数 HI_OTLP_ENDPOINT で指定する（必須）。
func Init(ctx context.Context, serviceName string) (func(context.Context) error, error) {
	exporterName := os.Getenv("HI_TRACES_EXPORTER")
	if exporterName == "" || exporterName == "none" {
		return func(context.Context) error { return nil }, nil
	}

	exp, err := newExporter(ctx, exporterName)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		// キューが満杯のときは Span を破棄せず、キューに追加できるまで待機する
		sdktrace.WithBatcher(exp, sdktrace.WithBlocking()),
		// 親のサンプリング判断にかかわらず、すべての Span を記録する
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(newResource(serviceName)),
	)
	otel.SetTracerProvider(tp)
	// これがないと OTel 内部のエラーが表示されないので、エクスポート失敗に気付けない
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		slog.Error("OpenTelemetry error", "err", err)
	}))
	return tp.Shutdown, nil
}

// newResource は service.name 属性を持つ Resource を生成する。
// resource.New は resource.Default() とマージされ host.name などの不要な属性が付くため使わない。
func newResource(serviceName string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
	)
}

// newExporter は名前に応じたエクスポーターを生成する。
// stdout は Span の JSON をコマンドの出力に混ぜないよう、標準エラー出力へ書く。
// otlp は HI_OTLP_ENDPOINT で指定した宛先へ OTLP/HTTP で送信する。
func newExporter(ctx context.Context, name string) (sdktrace.SpanExporter, error) {
	switch name {
	case "stdout":
		return stdouttrace.New(stdouttrace.WithWriter(os.Stderr))
	case "otlp":
		endpoint := os.Getenv("HI_OTLP_ENDPOINT")
		if endpoint == "" {
			return nil, errors.New("HI_OTLP_ENDPOINT is required when HI_TRACES_EXPORTER=otlp")
		}
		return otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	default:
		return nil, fmt.Errorf("unknown HI_TRACES_EXPORTER: %q (none|stdout|otlp)", name)
	}
}
