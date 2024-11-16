package transport

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"

	"kit-fiber-example/middlewares"
)

// Transport extension
type UppercaseRequest struct {
	S string `json:"string"`
}

func decodeUppercaseRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request UppercaseRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func (r UppercaseRequest) TraceAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("input", r.S),
	}
}

type UppercaseResponse struct {
	V   string `json:"result"`
	Err string `json:"error,omitempty"`
}

func decodeUppercaseResponse(_ context.Context, r *http.Response) (interface{}, error) {
	var response UppercaseResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func makeUppercaseEndpoint(svc StringService) middlewares.Endpoint[UppercaseRequest, UppercaseResponse] {
	return func(_ context.Context, req UppercaseRequest) (UppercaseResponse, error) {
		v, err := svc.Uppercase(req.S)
		if err != nil {
			return UppercaseResponse{v, err.Error()}, nil
		}
		return UppercaseResponse{v, ""}, nil
	}
}

// HandleUppercase is the Fiber handler for the uppercase endpoint
func (t *fiberTransport) HandleUppercase(c *fiber.Ctx) error {
	var req UppercaseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := t.Uppercase(c.Context(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
