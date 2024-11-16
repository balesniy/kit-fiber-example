package middlewares

import (
	"context"
	"time"

	"kit-fiber-example/metrics"
)

func metricsMiddleware[Req any, Res any](m *metrics.Metrics) Middleware[Req, Res] {
	return func(next Endpoint[Req, Res]) Endpoint[Req, Res] {
		return func(ctx context.Context, request Req) (Res, error) {
			defer func(begin time.Time) {
				m.RequestLatency.Observe(time.Since(begin).Seconds())
				m.RequestCount.Add(1)
			}(time.Now())

			result, err := next(ctx, request)
			if err != nil {
				m.ErrorCount.Add(1)
			}
			return result, err
		}
	}
}

func WithMetrics[Req any, Res any](m *metrics.Metrics, endpoint Endpoint[Req, Res]) Endpoint[Req, Res] {
	return metricsMiddleware[Req, Res](m)(endpoint)
}
