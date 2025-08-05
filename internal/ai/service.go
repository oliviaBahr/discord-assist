package ai

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

// Service handles AI interactions using Anthropic's Claude API
type Service struct {
	client             *anthropic.Client
	logger             *log.Logger
	model              string
	toolRegistry       *ToolRegistry
	defaultParams      anthropic.MessageNewParams
	sendDiscordMessage func(channelID, content string)
}

// NewService creates a new AI service
func NewService(apiKey, model string, logger *log.Logger, sendDiscordMessage func(channelID, content string)) (*Service, error) {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))

	// Create default parameters using proper SDK types
	defaultParams := anthropic.MessageNewParams{
		Model:       anthropic.Model(model),
		MaxTokens:   1000,
		System:      []anthropic.TextBlockParam{{Text: SystemPrompt}},
		Temperature: anthropic.Float(0.7),
		Tools:       GetTools(),
	}

	return &Service{
		client:             &client,
		logger:             logger,
		model:              model,
		toolRegistry:       GlobalToolRegistry,
		defaultParams:      defaultParams,
		sendDiscordMessage: sendDiscordMessage,
	}, nil
}

// createMessageParams creates MessageNewParams with default values and custom messages
func (s *Service) createMessageParams(messages []anthropic.MessageParam) anthropic.MessageNewParams {
	params := s.defaultParams
	params.Messages = messages
	return params
}

func (s *Service) buildConversationMessages(messages []*discordgo.Message) ([]anthropic.MessageParam, string, error) {
	if len(messages) == 0 {
		return nil, "", fmt.Errorf("no valid messages found in conversation context")
	}
	channelID := messages[0].ChannelID

	var conversationMessages []anthropic.MessageParam
	cleanedMessages := slices.DeleteFunc(messages, func(msg *discordgo.Message) bool {
		return msg.Content == ""
	})
	for _, msg := range slices.Backward(cleanedMessages) {
		var msgFunc = anthropic.NewUserMessage
		if msg.Author.Bot {
			msgFunc = anthropic.NewAssistantMessage
		}
		messageParam := msgFunc(anthropic.NewTextBlock(msg.Content))
		conversationMessages = append(conversationMessages, messageParam)
	}

	if len(conversationMessages) == 0 {
		return nil, "", fmt.Errorf("no valid messages found in conversation context")
	}

	return conversationMessages, channelID, nil
}

// GenerateResponse generates an AI response to a user message
func (s *Service) GenerateResponse(ctx context.Context, messages []*discordgo.Message) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("no messages provided")
	}

	conversationMessages, channelID, err := s.buildConversationMessages(messages)
	if err != nil {
		return "", err
	}

	for {
		s.logger.Info("generating response")
		resp, err := s.client.Messages.New(ctx, s.createMessageParams(conversationMessages))
		if err != nil || len(resp.Content) == 0 {
			return "", fmt.Errorf("failed to generate AI response: %w", err)
		}

		var toolUses []anthropic.ToolUseBlock
		var textBlocks []string
		for _, block := range resp.Content {
			switch block.Type {
			case "tool_use":
				toolUses = append(toolUses, block.AsToolUse())
			case "text":
				textBlocks = append(textBlocks, block.AsText().Text)
			}
		}

		s.sendDiscordMessage(channelID, strings.Join(textBlocks, "\n"))

		// Check if the response stopped due to tool use
		if resp.StopReason == "tool_use" {
			// First, add the assistant message with tool uses
			assistantBlocks := []anthropic.ContentBlockParamUnion{}
			for _, block := range resp.Content {
				assistantBlocks = append(assistantBlocks, block.ToParam())
			}
			assistantMessage := anthropic.NewAssistantMessage(assistantBlocks...)
			conversationMessages = append(conversationMessages, assistantMessage)

			// Then add tool results as user messages
			for _, toolUse := range toolUses {
				toolResultBlock, err := s.toolRegistry.ExecuteTool(toolUse.Name, toolUse.Input, toolUse.ID)
				if err != nil {
					return fmt.Sprintf("Sorry, I encountered an error while using a tool: %v", err), nil
				}
				toolResultMessage := anthropic.NewUserMessage(toolResultBlock)
				conversationMessages = append(conversationMessages, toolResultMessage)
			}
			// Continue the loop to send tool results and get the next response
			continue
		}

		// If the response stopped for any other reason (end_turn, max_tokens, etc.), we've already sent the text
		if len(textBlocks) > 0 {
			return "", nil // Text already sent to Discord
		}

		return "", fmt.Errorf("no text response in final message")
	}
}
