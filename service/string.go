package service

import (
	"context"
	"strings"
)

// Custom error types
type ServiceError struct {
	Code    int
	Message string
}

func (e ServiceError) Error() string {
	return e.Message
}

// stringService is a concrete implementation of StringService
type String struct {
	ClaudeClient *ClaudeClient
}

func (String) Uppercase(s string) (string, error) {
	return strings.ToUpper(s), nil
}

func (s *String) AskClaude(ctx context.Context, question string) (string, error) {
	// todo
	response, err = s.ClaudeClient.Ask(ctx, question)
	return "", nil
}
