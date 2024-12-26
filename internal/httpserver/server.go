package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gtuk/discordwebhook"

	"github.com/dwarvesf/teleport-discord-bot/internal/discord"
)

// Server represents an HTTP server for health checks and other utilities
type Server struct {
	server  *http.Server
	port    string
	mu      sync.Mutex
	discord *discord.Client
}

// NewServer creates a new HTTP server with a healthz endpoint
func NewServer(port string, discordClient *discord.Client) *Server {
	mux := http.NewServeMux()
	s := &Server{
		port:    port,
		discord: discordClient,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: mux,
		},
	}

	// Add healthz endpoint
	mux.HandleFunc("/healthz", s.healthzHandler)

	// Add event log endpoint
	mux.HandleFunc("/events.log", s.eventLogHandler)
	mux.HandleFunc("/session.log", s.sessionLogHandler)

	return s
}

// healthzHandler responds with a 200 OK status for health checks
func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// EventLogBody represents the structure of a Fluentd-like event log payload
type EventLogBody struct {
	AccessRequests   []string `json:"access_requests"`
	ClusterName      string   `json:"cluster_name"`
	Code             string   `json:"code"`
	DbName           string   `json:"db_name"`
	DbOrigin         string   `json:"db_origin"`
	DbProtocol       string   `json:"db_protocol"`
	DbQuery          string   `json:"db_query"`
	DbService        string   `json:"db_service"`
	DbType           string   `json:"db_type"`
	DbURI            string   `json:"db_uri"`
	DbUser           string   `json:"db_user"`
	Ei               int      `json:"ei"`
	Query            string   `json:"query"`
	Event            string   `json:"event"`
	PrivateKeyPolicy string   `json:"private_key_policy"`
	Sid              string   `json:"sid"`
	Success          bool     `json:"success"`
	UID              string   `json:"uid"`
	User             string   `json:"user"`
	UserKind         int      `json:"user_kind"`
}

func (s *Server) eventLogHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST method is accepted
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON body
	var body EventLogBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Print out the event log details in a structured format
	fmt.Printf("Fluentd Event Log Received:\n")
	fmt.Printf("User: %v, DbQuery: %v\n", body.User, body.DbQuery)

	query := ""
	switch body.Event {
	case "db.session.query":
		if body.DbQuery == "" {
			// Respond with success
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return
		}
		query = body.DbQuery
	case "db.session.postgres.statements.parse":
		if body.Query == "" {
			// Respond with success
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
			return
		}
		query = body.Query
	default:
		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		return
	}

	// Prepare Discord embed message
	if s.discord != nil {
		// Determine embed color based on success
		color := "15158332" // Red for failure
		if body.Success {
			color = "3066993" // Green for success
		}

		// Prepare fields for the embed
		fields := []discordwebhook.Field{
			{
				Name:   ptrString("User"),
				Value:  ptrString(body.User),
				Inline: ptrBool(true),
			},
		}

		// Add additional fields based on event type
		if body.DbName != "" {
			fields = append(fields,
				discordwebhook.Field{
					Name:   ptrString("Database User"),
					Value:  ptrString(body.DbUser),
					Inline: ptrBool(true),
				},
				discordwebhook.Field{
					Name:   ptrString("Database"),
					Value:  ptrString(body.DbName),
					Inline: ptrBool(true),
				},
				discordwebhook.Field{
					Name:   ptrString("Query"),
					Value:  ptrString(query),
					Inline: ptrBool(false),
				},
				discordwebhook.Field{
					Name:   ptrString("Success"),
					Value:  ptrString(fmt.Sprintf("%v", body.Success)),
					Inline: ptrBool(true),
				},
				discordwebhook.Field{
					Name:   ptrString("Time"),
					Value:  ptrString(time.Now().Format(time.RFC3339)),
					Inline: ptrBool(true),
				},
			)
		}

		// Create embed
		embed := discordwebhook.Embed{
			Title:       ptrString("Teleport Event Log"),
			Description: ptrString(fmt.Sprintf("Event details for %s", body.Event)),
			Color:       ptrString(color),
			Fields:      &fields,
		}

		// Prepare message
		message := discordwebhook.Message{
			Embeds: &[]discordwebhook.Embed{embed},
		}

		// Send webhook
		if err := s.discord.SendWebhookNotification(message); err != nil {
			fmt.Printf("Failed to send Discord webhook: %v\n", err)
		}
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event log received successfully"))
}

func (s *Server) sessionLogHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure only POST method is accepted
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	return
	// }

	// Decode JSON body
	var body map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Print out the event log details in a structured format
	fmt.Printf("Fluentd Session Log Received:\n")
	fmt.Printf("Session: %v\n", body)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Session log received successfully"))
}

// Helper functions for Discord webhook
func ptrString(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}

// Start starts the HTTP server in a separate goroutine
func (s *Server) Start() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("Starting HTTP server on port %v\n", s.port)
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *Server) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.server.Shutdown(ctx)
}
