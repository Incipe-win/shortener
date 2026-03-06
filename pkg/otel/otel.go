package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

const tracerName = "shortener-gateway"

// InitTracer 初始化 OpenTelemetry TracerProvider
// 返回 shutdown 函数，应在 main 中 defer 调用
func InitTracer(serviceName, endpoint string, samplerRatio float64) (func(), error) {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(endpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	var sampler sdktrace.Sampler
	if samplerRatio >= 1.0 {
		sampler = sdktrace.AlwaysSample()
	} else if samplerRatio <= 0 {
		sampler = sdktrace.NeverSample()
	} else {
		sampler = sdktrace.TraceIDRatioBased(samplerRatio)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)

	shutdown := func() {
		_ = tp.Shutdown(context.Background())
	}

	return shutdown, nil
}

// Tracer 返回全局 Tracer 实例
func Tracer() trace.Tracer {
	return otel.Tracer(tracerName)
}

// SpanAttrString 便捷函数：创建 string 类型的 span attribute
func SpanAttrString(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

// SpanAttrBool 便捷函数：创建 bool 类型的 span attribute
func SpanAttrBool(key string, value bool) attribute.KeyValue {
	return attribute.Bool(key, value)
}
