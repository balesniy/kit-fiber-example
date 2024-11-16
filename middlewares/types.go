package middlewares

import "context"

type Endpoint[Req any, Resp any] func(ctx context.Context, request Req) (response Resp, err error)
type Middleware[Req any, Res any] func(Endpoint[Req, Res]) Endpoint[Req, Res]

// Chain is a helper function for composing middlewares. Requests will
// traverse them in the order they're declared. That is, the first middleware
// is treated as the outermost middleware.
func Chain(outer Middleware[any, any], others ...Middleware[any, any]) Middleware[any, any] {
	return func(next Endpoint[any, any]) Endpoint[any, any] {
		for i := len(others) - 1; i >= 0; i-- { // reverse
			next = others[i](next)
		}
		return outer(next)
	}
}
