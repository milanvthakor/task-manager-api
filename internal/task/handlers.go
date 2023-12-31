package task

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/milanvthakor/task-manager-api/internal/models"
	"github.com/milanvthakor/task-manager-api/internal/validator"
	"github.com/milanvthakor/task-manager-api/pkg/config"
)

// taskData holds the task details.
type taskData struct {
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      models.TaskStatus `json:"status"`
}

// GetTasksHandler handles retrieval of a list of tasks associated with the authenticated user.
func GetTasksHandler(ctx *gin.Context, app *config.Application) {
	userID := ctx.MustGet("userID").(uint)

	// Retrieve tasks associated with the user from the database
	tasks, err := app.TaskRepository.ListTasksByUserID(userID)
	if err != nil {
		log.Printf("Warning: Failed to retrieve tasks: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve tasks"})
		return
	}

	ctx.JSON(http.StatusCreated, tasks)
}

// CreateTaskHandler handles the creation of a new task.
func CreateTaskHandler(ctx *gin.Context, app *config.Application) {
	userID := ctx.MustGet("userID").(uint)

	var td taskData
	if err := ctx.ShouldBindJSON(&td); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs"})
		return
	}

	// Validate inputs.
	if validator.IsBlank(td.Title) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid title. It must not be empty"})
		return
	}
	if !validator.IsValidTaskStatus(td.Status) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": `Invalid status. It can have one of the following values: "todo", "in progress", "done"`})
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Task created successfully",
		"task":    newTask,
	})
}

// GetTaskByIDHandler handles the retrieval of a task by ID only if it's associated with the authenticated user.
func GetTaskByIDHandler(ctx *gin.Context, app *config.Application) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.MustGet("taskID").(uint)

	// Retrieve the task from the database
	task, err := app.TaskRepository.GetTaskByID(uint(taskID))
	if err != nil {
		log.Printf("Warning: Failed to get task details from the database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	// Check if the task is associated with the authenticated user
	if task == nil || task.UserID != userID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// DeleteTaskByIDHandler handles the deletion of a task by ID only if it's associated with the authenticated user.
func DeleteTaskByIDHandler(ctx *gin.Context, app *config.Application) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.MustGet("taskID").(uint)

	// Delete the task from the database
	err := app.TaskRepository.DeleteTask(uint(taskID), userID)
	if err == sql.ErrNoRows {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	if err != nil {
		log.Printf("Warning: Failed to delete task from the database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// UpdateTaskByIDHandler handles the updating of a task by ID only if it's associated with the authentication user.
func UpdateTaskByIDHandler(ctx *gin.Context, app *config.Application) {
	userID := ctx.MustGet("userID").(uint)
	taskID := ctx.MustGet("taskID").(uint)

	// Retrieve the task from the database
	task, err := app.TaskRepository.GetTaskByID(uint(taskID))
	if err != nil {
		log.Printf("Warning: Failed to get task details from the database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve task"})
		return
	}

	// Check if the task is associated with the authenticated user
	if task == nil || task.UserID != userID {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Parse request body to get updated details
	var td taskData
	if err := ctx.ShouldBindJSON(&td); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs"})
		return
	}

	// Update only provided fields
	if !validator.IsBlank(td.Title) {
		task.Title = td.Title
	}
	if !validator.IsBlank(string(td.Status)) {
		if !validator.IsValidTaskStatus(td.Status) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": `Invalid status. It can have one of the following values: "todo", "in progress", "done"`})
			return
		}

		task.Status = td.Status
	}
	if !validator.IsBlank(td.Description) {
		task.Description = td.Description
	}

	// Update the task in the database
	updatedTask, err := app.TaskRepository.UpdateTask(task)
	if err != nil {
		log.Printf("Warning: Failed to update task: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Task updated successfully",
		"task":    updatedTask,
	})
}

// updateResult represents the result of each task update operation.
type updateResult struct {
	ID      uint   `json:"id"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// MarkTasksDoneHandler allows users to mark multiple tasks as "done".
func MarkTasksDoneHandler(ctx *gin.Context, app *config.Application) {
	userID := ctx.MustGet("userID").(uint)

	// Parse task IDs from the request body
	var taskIDs []uint
	if err := ctx.ShouldBindJSON(&taskIDs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID(s)"})
		return
	}

	// Create a channel to receive task update results
	updateResultChan := make(chan *updateResult, len(taskIDs))
	// Create a wait group to manage task update operation goroutines
	wg := sync.WaitGroup{}

	// Iterate over the task IDs and update them concurrently
	for _, taskID := range taskIDs {
		wg.Add(1)

		go func(taskID uint) {
			defer wg.Done()

			// Retrieve the task from the database
			task, err := app.TaskRepository.GetTaskByID(uint(taskID))
			if err != nil {
				log.Printf("Warning: Failed to get task details from the database: %v", err)
				updateResultChan <- &updateResult{ID: taskID, Error: "Failed to retrieve task"}
				return
			}

			// Check if the task is associated with the authenticated user
			if task == nil || task.UserID != userID {
				updateResultChan <- &updateResult{ID: taskID, Error: "Task not found"}
				return
			}

			// Update the task status to "done" and save it to the database
			task.Status = models.TaskStatusDone
			if _, err := app.TaskRepository.UpdateTask(task); err != nil {
				log.Printf("Warning: Failed to update task: %v", err)
				updateResultChan <- &updateResult{ID: taskID, Error: "Failed to update task"}
				return
			}

			updateResultChan <- &updateResult{ID: taskID, Message: "Task marked as done successfully"}
		}(taskID)
	}

	// Wait for all task update operation goroutines to finish
	wg.Wait()
	close(updateResultChan)

	// Collect results from goroutines
	var updatedResults []*updateResult
	for result := range updateResultChan {
		updatedResults = append(updatedResults, result)
	}

	ctx.JSON(http.StatusOK, updatedResults)
}
