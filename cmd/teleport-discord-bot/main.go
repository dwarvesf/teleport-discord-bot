package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dwarvesf/teleport-discord-bot/internal/config"
	"github.com/dwarvesf/teleport-discord-bot/internal/discord"
	"github.com/dwarvesf/teleport-discord-bot/internal/httpserver"
	repo "github.com/dwarvesf/teleport-discord-bot/internal/repository"
	"github.com/dwarvesf/teleport-discord-bot/internal/teleport"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	repo.ConnectDatabase()

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create Discord client
	discordClient := discord.NewClient(cfg)

	// Create HTTP server with Discord client
	httpServer := httpserver.NewServer(cfg.Port, discordClient)

	// Create Teleport plugin
	plugin, err := teleport.NewPlugin(cfg, discordClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Teleport plugin: %v\n", err)
		os.Exit(1)
	}

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown function
	gracefulShutdown := func() {
		fmt.Println("Initiating graceful shutdown...")

		// Cancel context to stop ongoing operations
		cancel()

		// Shutdown HTTP server
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "Error shutting down HTTP server: %v\n", err)
		}

		// Close Teleport plugin
		plugin.Close()

		fmt.Println("Teleport Discord bot shutdown complete")
		os.Exit(0)
	}

	// Run the plugin in a separate goroutine
	go func() {
		if err := plugin.Run(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Plugin error: %v\n", err)
			gracefulShutdown()
		}
	}()

	// Start the HTTP server
	if err := httpServer.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start HTTP server: %v\n", err)
		os.Exit(1)
	}

	// Wait for interrupt signal
	<-sigChan
	gracefulShutdown()
}
