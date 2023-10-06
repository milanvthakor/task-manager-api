package task

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/milanvthakor/task-manager-api/internal/models"
	"github.com/milanvthakor/task-manager-api/internal/validator"
	"github.com/milanvthakor/task-manager-api/pkg/config"
)

// TaskData holds the task details.
type TaskData struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
}

// validate validates the field values of the taskData struct.
func (t *TaskData) validate() error {
	if validator.IsBlank(t.Title) {
		return errors.New("Invalid title. It must not be empty")
	}

	if !validator.IsValidTaskStatus(t.Status) {
		return errors.New(`Invalid status. It can have one of the following values: "todo", "in progress", "done"`)
	}

	return nil
}

// CreateTaskHandler handles the creation of a new task.
func CreateTaskHandler(ctx *gin.Context, app *config.Application) {
	userID := ctx.MustGet("userID").(uint)

	var td TaskData
	if err := ctx.ShouldBindJSON(&td); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid inputs"})
		return
	}

	// Validate inputs.
	if err := td.validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Store task details in the database.
	task := &models.Task{
		Title:       td.Title,
		Description: td.Description,
		Status:      td.Status,
		UserID:      userID,
	}
	newTask, err := app.TaskRepository.CreateTask(task)
	if err != nil {
		log.Printf("Warning: Failed to create task: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Task created successfully",
		"task":    newTask,
	})
}

// GetTaskByIDHandler handles the retrieval of a task by ID only if it's associated with the authenticated user.
func GetTaskByIDHandler(ctx *gin.Context, app *config.Application) {
	userID := ctx.MustGet("userID").(uint)

	// Get the task ID from the URL parameters
	taskIDStr := ctx.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	// Retrieve the task from the database
	task, err := app.TaskRepository.GetTaskByID(uint(taskID))
	if err != nil {
		log.Printf("Warning: Failed to get task details from the database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to retrieve task"})
		return
	}

	// Check if the task is associated with the authenticated user
	if task == nil || task.UserID != userID {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}
