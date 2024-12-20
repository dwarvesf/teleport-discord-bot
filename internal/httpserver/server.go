package httpserver

import (
	"context"
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

	return s
}

// healthzHandler responds with a 200 OK status for health checks
func (s *Server) healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
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
