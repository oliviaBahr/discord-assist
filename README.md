# Discord Assist

A Discord bot with a menu bar interface for easy control.

## Features

- Discord bot with AI-powered responses
- System tray menu bar for easy control
- Start/Stop bot functionality
- Quit application option

## Building

```bash
go build
```

## Running

```bash
./discord-assist
```

The application will appear as a menu bar icon. You can:

- **Start Bot**: Start the Discord bot
- **Stop Bot**: Stop the Discord bot  
- **Quit**: Exit the application

## Configuration

Make sure to set up your `.env` file with the required Discord and Anthropic API keys:

```
DISCORD_TOKEN=your_discord_bot_token
ANTHROPIC_API_KEY=your_anthropic_api_key
```

## Menu Bar Features

The menu bar provides easy access to control the Discord bot:

- The bot starts automatically when the application launches
- Use the menu to start/stop the bot as needed
- The application runs in the background with a system tray icon
- Click "Quit" to completely exit the application

## Adding an Icon

To add a custom icon to the menu bar:

1. Create a PNG icon file at `assets/icon.png`
2. The icon should be 16x16 or 32x32 pixels for best results
3. Rebuild the application with `go build`

The menu bar will automatically use the custom icon when available. 