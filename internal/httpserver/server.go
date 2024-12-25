package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

// Server represents an HTTP server for health checks and other utilities
type Server struct {
	server *http.Server
	port   string
	mu     sync.Mutex
}

// NewServer creates a new HTTP server with a healthz endpoint
func NewServer(port string) *Server {
	mux := http.NewServeMux()
	s := &Server{
		port: port,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%v", port),
			Handler: mux,
		},
	}

	// Add healthz endpoint
	mux.HandleFunc("/healthz", s.healthzHandler)

	// Add event log endpoint
	mux.HandleFunc("/events.log", s.eventLogHandler)

	return s
}

// healthzHandler responds with a 200 OK status for health checks
func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// EventLogBody represents the structure of a Fluentd-like event log payload
type EventLogBody struct {
	Tag    string                 `json:"tag"`
	Time   float64                `json:"time"`
	Record map[string]interface{} `json:"record"`
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
	fmt.Printf("Tag: %s\n", body.Tag)
	fmt.Printf("Timestamp: %f\n", body.Time)
	fmt.Println("Record:")
	for key, value := range body.Record {
		fmt.Printf("  %s: %v\n", key, value)
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Event log received successfully"))
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
