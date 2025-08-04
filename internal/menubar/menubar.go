package menubar

import (
	"context"
	"log/slog"
	"os"

	"github.com/getlantern/systray"
)

// MenuBar represents the system tray menu bar
type MenuBar struct {
	logger *slog.Logger
	bot    BotController
}

// BotController defines the interface for controlling the bot
type BotController interface {
	Start(ctx context.Context) error
	Stop() error
	IsRunning() bool
}

// New creates a new menu bar instance
func New(bot BotController, logger *slog.Logger) *MenuBar {
	return &MenuBar{
		logger: logger,
		bot:    bot,
	}
}

// Run starts the menu bar application
func (m *MenuBar) Run() {
	systray.Run(m.onReady, m.onExit)
}

// onReady is called when the systray is ready
func (m *MenuBar) onReady() {
	// Load icon from embedded assets
	iconBytes, err := m.loadIcon()
	if err != nil || len(iconBytes) == 0 {
		m.logger.Error("failed to load icon", "error", err)
		// Use a default icon or continue without icon
		iconBytes = []byte{}
	}

	systray.SetIcon(iconBytes)
	systray.SetTooltip("Discord Assist Bot")

	// Add menu items
	mQuit := systray.AddMenuItem("Quit Discord Assist", "Quit the application")
	systray.AddSeparator()
	mStart := systray.AddMenuItem("Start Bot", "Start the Discord bot")
	mStop := systray.AddMenuItem("Stop Bot", "Stop the Discord bot")

	// Start the bot automatically
	go func() {
		ctx := context.Background()
		if err := m.bot.Start(ctx); err != nil {
			m.logger.Error("failed to start bot", "error", err)
		}
	}()

	// Handle menu item clicks
	go func() {
		for {
			select {
			case <-mQuit.ClickedCh:
				m.logger.Info("quit requested from menu")
				systray.Quit()
				return
			case <-mStart.ClickedCh:
				if !m.bot.IsRunning() {
					go func() {
						ctx := context.Background()
						if err := m.bot.Start(ctx); err != nil {
							m.logger.Error("failed to start bot", "error", err)
						}
					}()
				}
			case <-mStop.ClickedCh:
				if m.bot.IsRunning() {
					if err := m.bot.Stop(); err != nil {
						m.logger.Error("failed to stop bot", "error", err)
					}
				}
			}
		}
	}()
}

// onExit is called when the systray is exiting
func (m *MenuBar) onExit() {
	m.logger.Info("menu bar exiting")
	// Stop the bot when the menu bar exits
	if m.bot.IsRunning() {
		if err := m.bot.Stop(); err != nil {
			m.logger.Error("failed to stop bot on exit", "error", err)
		}
	}
}

// loadIcon loads the icon from assets
func (m *MenuBar) loadIcon() ([]byte, error) {
	// Read PNG file directly
	iconBytes, err := os.ReadFile("assets/icon.png")
	if err != nil {
		m.logger.Error("failed to load icon", "error", err)
		return []byte{}, nil
	}
	return iconBytes, nil
}
