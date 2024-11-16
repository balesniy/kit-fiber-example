package middlewares

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type TracingRequest interface {
	TraceAttributes() []attribute.KeyValue
}

func tracingMiddleware[Req any, Res any](tracer trace.Tracer) Middleware[Req, Res] {
	return func(next Endpoint[Req, Res]) Endpoint[Req, Res] {
		return func(ctx context.Context, request Req) (Res, error) {
			spanCtx, span := tracer.Start(ctx, "endpoint")
			defer span.End()

			if tr, ok := any(request).(TracingRequest); ok {
				span.SetAttributes(tr.TraceAttributes()...)
			}

			result, err := next(spanCtx, request)
			if err != nil {
				span.RecordError(err)
			}

			return result, err
		}
	}
}

func WithTracing[Req any, Res any](tracer trace.Tracer, endpoint Endpoint[Req, Res]) Endpoint[Req, Res] {
	return tracingMiddleware[Req, Res](tracer)(endpoint)
}
