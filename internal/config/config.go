package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	ProxyAddr         string
	DiscordWebhookURL string
	WatcherList       string
	AuthPemPath       string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		ProxyAddr:         os.Getenv("PROXY_ADDR"),
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		WatcherList:       os.Getenv("WATCHER_LIST"),
		AuthPemPath:       os.Getenv("AUTH_PEM_PATH"),
	}

	// Validate required configuration
	if cfg.ProxyAddr == "" {
		return nil, fmt.Errorf("PROXY_ADDR is required")
	}
	if cfg.DiscordWebhookURL == "" {
		return nil, fmt.Errorf("DISCORD_WEBHOOK_URL is required")
	}
	if cfg.AuthPemPath == "" {
		return nil, fmt.Errorf("AUTH_PEM_PATH is required")
	}

	return cfg, nil
}
