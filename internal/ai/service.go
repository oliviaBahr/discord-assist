package ai

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// Service handles AI interactions using Anthropic's Claude API
type Service struct {
	client *anthropic.Client
	logger *slog.Logger
	model  string
}

// NewService creates a new AI service
func NewService(apiKey, model string, logger *slog.Logger) (*Service, error) {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	return &Service{
		client: &client,
		logger: logger,
		model:  model,
	}, nil
}

// GenerateResponse generates an AI response to a user message
func (s *Service) GenerateResponse(ctx context.Context, userMessage, username string) (string, error) {
	// Create a system prompt that defines the bot's behavior
	systemPrompt := `You are a helpful Discord bot assistant. You should:
- Be friendly and conversational
- Keep responses concise but helpful
- Be appropriate for a Discord chat environment
- Respond naturally to questions and statements
- Use emojis occasionally to make responses more engaging
- Don't be overly formal unless the user is asking for something technical`

	// Create the message for Claude
	message := anthropic.NewUserMessage(
		anthropic.NewTextBlock(fmt.Sprintf("User %s says: %s", username, userMessage)),
	)

	// Make the API call
	resp, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:       anthropic.Model(s.model),
		MaxTokens:   500,
		Messages:    []anthropic.MessageParam{message},
		System:      []anthropic.TextBlockParam{{Text: systemPrompt}},
		Temperature: anthropic.Float(0.7),
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate AI response: %w", err)
	}

	if len(resp.Content) == 0 {
		return "I'm sorry, I couldn't generate a response right now.", nil
	}

	// Extract the text response
	response := resp.Content[0].Text
	s.logger.Debug("generated AI response",
		"user", username,
		"input_length", len(userMessage),
		"response_length", len(response))

	return response, nil
}
