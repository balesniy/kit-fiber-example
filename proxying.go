package main

import (
	"context"
	"net/http"
	"net/url"

	"kit-fiber-example/middlewares"
	"kit-fiber-example/transport"
)

// proxymw implements StringService, forwarding Uppercase requests to the
// provided endpoint, and serving all other (i.e. Count) requests via the
// next StringService.
type proxymw struct {
	next      transport.StringService                                                       // Serve all requests via base service...
	askClaude middlewares.Endpoint[transport.AskClaudeRequest, transport.AskClaudeResponse] // ...except Claude, which gets served by this endpoint
}

// We’ve got exactly the same endpoint, but we’ll use it to invoke, rather than serve, a request.
func (mw proxymw) AskClaude(ctx context.Context, question string) (string, error) {
	response, err := mw.askClaude(ctx, transport.AskClaudeRequest{Question: question})
	if err != nil {
		return "", err
	}
	return response.Answer, nil
}

func (mw proxymw) Uppercase(s string) (string, error) {
	return mw.next.Uppercase(s)
}

type ServiceMiddleware func(transport.StringService) transport.StringService

func proxyingMiddleware(proxyURL string) ServiceMiddleware {
	return func(next transport.StringService) transport.StringService {
		return proxymw{next, makeClaudeEndpoint(proxyURL)}
	}
}

// assume JSON over HTTP
func makeClaudeEndpoint(proxyURL string) middlewares.Endpoint[transport.AskClaudeRequest, transport.AskClaudeResponse] {
	return func(ctx context.Context, request transport.AskClaudeRequest) (transport.AskClaudeResponse, error) {
		ctx, cancel := context.WithCancel(ctx)
		var (
			resp *http.Response
			err  error
		)
		c := &Client{
			client: http.DefaultClient,
			req:    makeCreateRequestFunc("POST", url.Parse(proxyURL), encodeClaudeRequest),
			dec:    decodeClaudeResponse,
		}

		req, err := c.req(ctx, request)
		if err != nil {
			cancel()
			return nil, err
		}
		resp, err = c.client.Do(req.WithContext(ctx))
		if err != nil {
			cancel()
			return nil, err
		}
		defer resp.Body.Close()
		defer cancel()
		response, err := c.dec(ctx, resp)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}

// HTTPClient is an interface that models *http.Client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Client wraps a URL and provides a method that implements endpoint.Endpoint.
type Client struct {
	client HTTPClient
	req    CreateRequestFunc
	dec    DecodeResponseFunc
}

type EncodeRequestFunc func(context.Context, *http.Request, interface{}) error
type CreateRequestFunc func(context.Context, interface{}) (*http.Request, error)
type DecodeResponseFunc func(context.Context, *http.Response) (interface{}, error)

func makeCreateRequestFunc(method string, target *url.URL, enc EncodeRequestFunc) CreateRequestFunc {
	return func(ctx context.Context, request interface{}) (*http.Request, error) {
		req, err := http.NewRequest(method, target.String(), nil)
		if err != nil {
			return nil, err
		}

		if err = enc(ctx, req, request); err != nil {
			return nil, err
		}

		return req, nil
	}
}

type ClientOption func(*Client)

func NewExplicitClient(req CreateRequestFunc, dec DecodeResponseFunc, options ...ClientOption) *Client {
	c := &Client{
		client: http.DefaultClient,
		req:    req,
		dec:    dec,
	}
	for _, option := range options {
		option(c)
	}
	return c
}
