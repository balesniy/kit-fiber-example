package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"kit-fiber-example/middlewares"
)

// Transport extension
type AskClaudeRequest struct {
	Question string `json:"question"`
}

type AskClaudeResponse struct {
	Answer string `json:"answer"`
	Error  string `json:"error,omitempty"`
}

func encodeClaudeRequest(_ context.Context, r *http.Request, request any) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(request); err != nil {
		return err
	}
	r.Body = io.NopCloser(&buf)
	return nil
}

func decodeClaudeResponse(_ context.Context, r *http.Response) (any, error) {
	var response AskClaudeResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response, nil
}

func makeAskClaudeEndpoint(svc StringService) middlewares.Endpoint[AskClaudeRequest, AskClaudeResponse] {
	return func(ctx context.Context, req AskClaudeRequest) (AskClaudeResponse, error) {
		answer, err := svc.AskClaude(ctx, req.Question)
		if err != nil {
			return AskClaudeResponse{"", err.Error()}, nil
		}
		return AskClaudeResponse{answer, ""}, nil
	}
}

func (t *fiberTransport) HandleAskClaude(c *fiber.Ctx) error {
	var req AskClaudeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	response, err := t.AskClaude(context.TODO(), req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(response)
}
