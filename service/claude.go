package service

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"kit-fiber-example/config"
)

// Claude API structures
type ClaudeRequest struct {
	Model     string    `json:"model"`
	Messages  []Message `json:"messages"`
	MaxTokens int       `json:"max_tokens,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeResponse struct {
	ID           string    `json:"id"`
	Model        string    `json:"model"`
	Role         string    `json:"role"`
	Content      []Content `json:"content"`
	Usage        Usage     `json:"usage"`
	CreatedAt    string    `json:"created_at"`
	CompletedAt  string    `json:"completed_at"`
	StopReason   string    `json:"stop_reason"`
	StopSequence string    `json:"stop_sequence"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type ClaudeClient struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
	model      string
}

func NewClaudeClient(cfg *config.Config) *ClaudeClient {
	return &ClaudeClient{
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Claude.Timeout) * time.Second,
		},
		apiKey:  cfg.Claude.APIKey,
		baseURL: cfg.Claude.BaseURL,
		model:   cfg.Claude.Model,
	}
}

func (c *ClaudeClient) Ask(ctx context.Context, question string) (string, error) {
	request := ClaudeRequest{
		Model: c.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: question,
			},
		},
	}

	postBody, err := json.Marshal(request)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL, bytes.NewBuffer(postBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil
	}

	var response ClaudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", err
	}

	return response.Content[0].Text, nil
}
