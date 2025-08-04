package bot

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"discord-assist/internal/ai"
	"discord-assist/internal/config"
	"discord-assist/internal/discord"
	"discord-assist/pkg/logger"
)

// Bot represents the main bot instance
type Bot struct {
	config  *config.Config
	client  *discord.Client
	ai      *ai.Service
	logger  *slog.Logger
	running bool
}

// New creates a new bot instance
func New(cfg *config.Config) (*Bot, error) {
	// Set up logger
	log := logger.Setup(cfg.Logging.Level, cfg.Logging.Format)

	// Create Discord client
	client, err := discord.NewClient(cfg.Discord.Token, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord client: %w", err)
	}

	// Create AI service
	aiService, err := ai.NewService(cfg.Anthropic.APIKey, cfg.Anthropic.Model, log)
	if err != nil {
		return nil, fmt.Errorf("failed to create AI service: %w", err)
	}

	bot := &Bot{
		config: cfg,
		client: client,
		ai:     aiService,
		logger: log,
	}

	// Set up event handlers
	bot.setupEventHandlers()

	return bot, nil
}

// Start starts the bot
func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("starting Discord bot...")
	b.running = true

	// Connect to Discord
	if err := b.client.Connect(ctx); err != nil {
		b.running = false
		return fmt.Errorf("failed to connect to Discord: %w", err)
	}

	// Set bot activity
	if err := b.client.SetActivity(b.config.Bot.ActivityType, b.config.Bot.Activity); err != nil {
		b.logger.Warn("failed to set bot activity", "error", err)
	}

	b.logger.Info("bot started successfully")

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan
	b.logger.Info("received shutdown signal, stopping bot...")

	return b.Stop()
}

// Stop stops the bot gracefully
func (b *Bot) Stop() error {
	b.logger.Info("stopping bot...")
	b.running = false

	if err := b.client.Close(); err != nil {
		b.logger.Error("error closing Discord connection", "error", err)
		return err
	}

	b.logger.Info("bot stopped successfully")
	return nil
}

// IsRunning returns whether the bot is currently running
func (b *Bot) IsRunning() bool {
	return b.running
}

// Logger returns the bot's logger
func (b *Bot) Logger() *slog.Logger {
	return b.logger
}

// setupEventHandlers sets up all event handlers
func (b *Bot) setupEventHandlers() {
	session := b.client.Session()

	// Message create event
	session.AddHandler(b.handleMessageCreate)
}

// handleMessageCreate handles message create events
func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore messages from bots
	if m.Author.Bot {
		return
	}

	b.logger.Debug("received message",
		"author", m.Author.Username,
		"content", m.Content,
		"channel", m.ChannelID,
	)

	// Generate AI response
	response, err := b.ai.GenerateResponse(context.Background(), m.Content, m.Author.Username)
	if err != nil {
		b.logger.Error("failed to generate AI response", "error", err)
		response = "I'm sorry, I'm having trouble processing your message right now. 😅"
	}

	// Send the response
	if err := b.client.SendMessage(m.ChannelID, response); err != nil {
		b.logger.Error("failed to send AI response", "error", err)
	}
}
