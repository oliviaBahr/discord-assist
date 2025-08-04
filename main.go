package main

import (
	"log"

	"discord-assist/internal/bot"
	"discord-assist/internal/config"
	"discord-assist/internal/menubar"
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

	// Create menu bar
	menuBar := menubar.New(b, b.Logger())

	// Run the menu bar (this will start the bot automatically)
	menuBar.Run()
}
