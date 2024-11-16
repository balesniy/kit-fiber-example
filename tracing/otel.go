package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"kit-fiber-example/config"
)

func InitOtel(cfg *config.Config) (*sdktrace.TracerProvider, error) {
	ctx := context.Background()

	conn, err := grpc.NewClient(cfg.Telemetry.CollectorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	//exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL("http://jaeger:4318"))
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			//semconv.ServiceNameKey.String(cfg.Telemetry.ServiceName),
			semconv.ServiceName(cfg.Telemetry.ServiceName),
			semconv.URLFull("jaeger"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.Telemetry.SamplingRatio)),
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(r),
	)

	otel.SetTracerProvider(tp)
	return tp, nil
}
