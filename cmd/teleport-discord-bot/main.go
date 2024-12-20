package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dwarvesf/teleport-discord-bot/internal/config"
	"github.com/dwarvesf/teleport-discord-bot/internal/discord"
	"github.com/dwarvesf/teleport-discord-bot/internal/httpserver"
	"github.com/dwarvesf/teleport-discord-bot/internal/teleport"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Create HTTP server
	httpServer := httpserver.NewServer(cfg.Port)
	defer func() {
		if shutdownErr := httpServer.Shutdown(context.Background()); shutdownErr != nil {
			fmt.Fprintf(os.Stderr, "Error shutting down HTTP server: %v\n", shutdownErr)
		}
	}()

	// Create Discord client
	discordClient := discord.NewClient(cfg)

	// Create Teleport plugin
	plugin, err := teleport.NewPlugin(cfg, discordClient)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Teleport plugin: %v\n", err)
		os.Exit(1)
	}
	defer plugin.Close()

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Run the plugin in a separate goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- plugin.Run(ctx)
	}()

	// Start the HTTP server
	if err := httpServer.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start HTTP server: %v\n", err)
		os.Exit(1)
	}

	// Wait for either a signal or an error
	select {
	case sig := <-sigChan:
		fmt.Printf("Received signal %v, shutting down...\n", sig)
		cancel()
	case err := <-errChan:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Plugin error: %v\n", err)
			os.Exit(1)
		}
	}

	// Wait for the plugin to finish
	<-errChan
	fmt.Println("Teleport Discord bot shutdown complete")
}
