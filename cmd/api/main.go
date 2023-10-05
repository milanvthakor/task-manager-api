package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/milanvthakor/task-manager-api/pkg/api"
	"github.com/milanvthakor/task-manager-api/pkg/config"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could load .env file %v\n", err)
	}

	// Initialize the configuration.
	cfg := config.New()

	// Initialize the Gin router.
	r := gin.Default()

	// Simple health check endpoint.
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Server is up and running")
	})

	// Start the server on the specified port.
	server := api.NewServer(r, cfg.ServerPort)
	server.Start()
}