package middlewares

import (
	"context"
	"fmt"
	"time"
)

func LoggingMiddleware[Req any, Res any](next Endpoint[Req, Res]) Endpoint[Req, Res] {
	return func(ctx context.Context, request Req) (Res, error) {
		start := time.Now()
		result, err := next(ctx, request)
		duration := time.Since(start)

		// Log the request
		fmt.Printf(
			"method=%s path=%s duration=%s err=%v\n",
			ctx.Value("method").(string),
			ctx.Value("path").(string),
			duration,
			err,
		)

		return result, err
	}
}

/*func anotherLoggingMiddleware[Req any, Res any](logger *log.Logger) Middleware[Req, Res] {
	return func(next Endpoint[Req, Res]) Endpoint[Req, Res] {
		return func(ctx context.Context, request Req) (Res, error) {
			logger.Printf("calling endpoint") // log.With(logger, "method", "Uppercase")
			defer logger.Printf("called endpoint")
			return next(ctx, request)
		}
	}
}*/
