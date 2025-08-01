package main

import (
	"context"
	"log"
	"os"

	"discord-assist/internal/bot"
	"discord-assist/internal/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create bot instance
	b, err := bot.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	// Create context
	ctx := context.Background()

	// Start the bot
	if err := b.Start(ctx); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	os.Exit(0)
} 