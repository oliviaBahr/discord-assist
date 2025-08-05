package discord

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
)

// Client wraps the Discord session and provides additional functionality
type Client struct {
	session *discordgo.Session
	logger  *log.Logger
}

// NewClient creates a new Discord client
func NewClient(token string, logger *log.Logger) (*Client, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	client := &Client{
		session: session,
		logger:  logger,
	}

	// Set up event handlers
	client.setupEventHandlers()

	return client, nil
}

// Connect establishes a connection to Discord
func (c *Client) Connect(ctx context.Context) error {
	c.logger.Info("connecting to Discord...")

	if err := c.session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord connection: %w", err)
	}

	c.logger.Info("successfully connected to Discord")
	return nil
}

// Close closes the Discord connection
func (c *Client) Close() error {
	c.logger.Info("closing Discord connection...")
	return c.session.Close()
}

// Session returns the underlying Discord session
func (c *Client) Session() *discordgo.Session {
	return c.session
}

// setupEventHandlers sets up the basic event handlers
func (c *Client) setupEventHandlers() {
	// Ready event
	c.session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		c.logger.Info("bot is ready",
			"user", s.State.User.Username,
			"guilds", len(s.State.Guilds),
		)
	})

	// Disconnect event
	c.session.AddHandler(func(s *discordgo.Session, d *discordgo.Disconnect) {
		c.logger.Warn("disconnected from Discord")
	})

	// Reconnect event
	c.session.AddHandler(func(s *discordgo.Session, r *discordgo.Connect) {
		c.logger.Info("reconnected to Discord")
	})
}

// SetActivity sets the bot's activity
func (c *Client) SetActivity(activityType, activity string) error {
	var discordActivityType discordgo.ActivityType

	switch activityType {
	case "Playing":
		discordActivityType = discordgo.ActivityTypeGame
	case "Streaming":
		discordActivityType = discordgo.ActivityTypeStreaming
	case "Listening":
		discordActivityType = discordgo.ActivityTypeListening
	case "Watching":
		discordActivityType = discordgo.ActivityTypeWatching
	default:
		discordActivityType = discordgo.ActivityTypeGame
	}

	return c.session.UpdateGameStatus(int(discordActivityType), activity)
}

// SendMessage sends a message to a channel
func (c *Client) SendMessage(channelID, content string) error {
	_, err := c.session.ChannelMessageSend(channelID, content)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

// SendEmbed sends an embed message to a channel
func (c *Client) SendEmbed(channelID string, embed *discordgo.MessageEmbed) error {
	_, err := c.session.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		return fmt.Errorf("failed to send embed: %w", err)
	}
	return nil
}

// GetRecentMessages fetches the last N messages from a channel
func (c *Client) GetRecentMessages(channelID string, limit int) ([]*discordgo.Message, error) {
	messages, err := c.session.ChannelMessages(channelID, limit, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recent messages: %w", err)
	}
	return messages, nil
}
