package task

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ExtractTaskIDMiddleware extract the task ID from URL parameters.
func ExtractTaskIDMiddleware(ctx *gin.Context) {
	taskIDStr := ctx.Param("id")
	taskID, err := strconv.ParseUint(taskIDStr, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid task ID"})
		return
	}

	// Store the task ID in the context
	ctx.Set("taskID", uint(taskID))
	ctx.Next()
}
