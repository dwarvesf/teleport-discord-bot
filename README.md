# Teleport Discord Access Request Notifier

## Overview

A Teleport plugin that sends real-time Discord webhook notifications for access requests, providing instant visibility into access management events.

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

Create a `.env` file with the following environment variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `PROXY_ADDR` | Teleport proxy address | `teleport.example.com:443` |
| `DISCORD_WEBHOOK_URL` | Discord webhook URL for notifications | `https://discord.com/api/webhooks/...` |
| `WATCHER_LIST` | Command to list watchers | `tctl` |
| `AUTH_PEM` | Base64-encoded Teleport authentication credentials | `LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVkt...` |

## Installation

1. Clone the repository
2. Copy `.env.example` to `.env` and configure your settings
3. Build the application:
   ```bash
   make build
   ```

## Running the Application

```bash
make run
```

## Development

### Setup Development Environment

```bash
make setup
```

### Available Make Targets

- `make build`: Compile the application
- `make test`: Run tests
- `make lint`: Run code linters
- `make run`: Build and run the application
- `make help`: Show all available commands

## Project Structure

```
teleport-discord-bot/
├── cmd/
│   └── teleport-discord-bot/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── discord/              # Discord webhook client
│   ├── teleport/             # Teleport event monitoring
├── go.mod
└── README.md
```

## Contributing

Contributions are welcome! Please follow these steps:
1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

Apache License 2.0
