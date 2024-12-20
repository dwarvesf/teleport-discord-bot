# Teleport Discord Access Request Notifier

## Overview

This is a Teleport Access Request plugin that sends notifications to a Discord webhook when access requests are created, approved, or denied.

## Features

- Real-time monitoring of Teleport access requests
- Discord webhook notifications for:
  - New access requests
  - Approved access requests
  - Denied access requests

## Prerequisites

- Go 1.21+
- Teleport cluster
- Discord webhook URL
- Teleport authentication credentials (auth.pem)

## Configuration

Create a `.env` file or set the following environment variables:

- `PROXY_ADDR`: Teleport proxy address (e.g., `teleport.example.com:443`)
- `DISCORD_WEBHOOK_URL`: Discord webhook URL for notifications
- `WATCHER_LIST`: Command to list watchers (used in notification messages)
- `AUTH_PEM_PATH`: Path to Teleport authentication credentials file

## Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and fill in your configuration
3. Build the application:
   ```bash
   go build -o github.com/dwarvesf/teleport-discord-bot ./cmd/github.com/dwarvesf/teleport-discord-bot
   ```

## Running the Application

```bash
./github.com/dwarvesf/teleport-discord-bot
```

## Project Structure

```
github.com/dwarvesf/teleport-discord-bot/
├── cmd/
│   └── github.com/dwarvesf/teleport-discord-bot/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── discord/
│   │   └── client.go        # Discord webhook client
│   ├── teleport/
│   │   └── plugin.go        # Teleport event monitoring
│   └── watcher/
│       └── watcher.go       # Event watching utilities
├── go.mod
└── README.md
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

Apache License 2.0
