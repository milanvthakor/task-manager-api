package task

import (
	"errors"
	"log"
	"net/http"

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
