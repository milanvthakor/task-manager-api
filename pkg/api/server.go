package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Server represents the API server.
type Server struct {
	router *gin.Engine
	port   string
}

// NewServer creates a new instance of the API server.
func NewServer(router *gin.Engine, port string) *Server {
	return &Server{
		router: router,
		port:   port,
	}
}

// Start starts the API server.
func (s *Server) Start() {
	addr := ":" + s.port
	log.Printf("Server is running on port %s...\n", s.port)
	if err := http.ListenAndServe(addr, s.router); err != nil {
		log.Fatalf("failed to start the server: %v", err)
	}
}
