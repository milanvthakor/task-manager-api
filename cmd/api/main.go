package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/milanvthakor/task-manager-api/internal/auth"
	"github.com/milanvthakor/task-manager-api/internal/database"
	"github.com/milanvthakor/task-manager-api/internal/models"
	"github.com/milanvthakor/task-manager-api/internal/task"
	"github.com/milanvthakor/task-manager-api/internal/utils"
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
		Config:         cfg,
		UserRepository: models.NewUserRepository(db),
		TaskRepository: models.NewTaskRepository(db),
	}

	// Initialize the Gin router.
	r := gin.Default()

	// Set up API routes
	apiRoutes := r.Group("/api")
	apiRoutes.POST("/register", utils.InjectApp(app, auth.ValidateInputMiddleware), utils.InjectApp(app, auth.RegisterHandler))
	apiRoutes.POST("/login", utils.InjectApp(app, auth.ValidateInputMiddleware), utils.InjectApp(app, auth.LoginHandler))
	// Set up Task API routes
	taskApiRoutes := apiRoutes.Group("/tasks")
	taskApiRoutes.POST("/", utils.InjectApp(app, auth.AuthenticateMiddleware), utils.InjectApp(app, task.CreateTaskHandler))
	taskApiRoutes.GET("/:id", utils.InjectApp(app, auth.AuthenticateMiddleware), task.ExtractTaskIDMiddleware, utils.InjectApp(app, task.GetTaskByIDHandler))
	taskApiRoutes.DELETE("/:id", utils.InjectApp(app, auth.AuthenticateMiddleware), task.ExtractTaskIDMiddleware, utils.InjectApp(app, task.DeleteTaskByIDHandler))
	taskApiRoutes.PUT("/:id", utils.InjectApp(app, auth.AuthenticateMiddleware), task.ExtractTaskIDMiddleware, utils.InjectApp(app, task.UpdateTaskByIDHandler))

	// Simple health check endpoint.
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "Server is up and running")
	})

	// Start the server on the specified port.
	server := api.NewServer(r, cfg.ServerPort)
	server.Start()
}
