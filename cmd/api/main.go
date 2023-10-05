package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/milanvthakor/task-manager-api/internal/auth"
	"github.com/milanvthakor/task-manager-api/internal/database"
	"github.com/milanvthakor/task-manager-api/internal/models"
	"github.com/milanvthakor/task-manager-api/pkg/api"
	"github.com/milanvthakor/task-manager-api/pkg/config"
)

func main() {
	// Load environment variables from the .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Couldn't load .env file %v\n", err)
	}

	// Initialize the configuration.
	cfg := config.New()

	// Initialize the database.
	db, err := database.Init(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to initialize the database: %v", err)
	}
	defer database.Close()

	// Initialize the new instance of the Application struct containing dependencies
	app := &config.Application{
		UserRepository: *models.NewUserRepository(db),
	}

	// Initialize the Gin router.
	r := gin.Default()

	// Set up API routes
	apiRoutes := r.Group("/api")
	apiRoutes.POST("/register", auth.RegisterHandler(app))

	// Simple health check endpoint.
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Server is up and running")
	})

	// Start the server on the specified port.
	server := api.NewServer(r, cfg.ServerPort)
	server.Start()
}
