package api

import (
	"net/http"

	"github.com/Mark-Life/smart-webp-resize/pkg/config"
)

// Server represents the HTTP server that handles image processing requests
type Server struct {
	config *config.Config
	router http.ServeMux
}

// NewServer creates a new server instance and sets up routes
func NewServer(cfg *config.Config) *Server {
	s := &Server{
		config: cfg,
		router: http.ServeMux{},
	}

	s.setupRoutes()
	return s
}

// setupRoutes configures all the routes for the server
func (s *Server) setupRoutes() {
	s.router.HandleFunc("/health", s.handleHealth)
	s.router.HandleFunc("/resize", s.handleResize)
}

// ServeHTTP makes the server implement the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// handleHealth handles health check requests
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
}

// handleResize handles image resize requests
func (s *Server) handleResize(w http.ResponseWriter, r *http.Request) {
	// Not implemented yet
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte(`{"error":"Not implemented yet"}`))
} 