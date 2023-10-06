package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milanvthakor/task-manager-api/pkg/config"
)

// Middleware for authentication input validation
func ValidateInput(ctx *gin.Context, app *config.Application) {
	var ud userData
	if err := ctx.ShouldBindJSON(&ud); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid inputs"})
		return
	}

	// Validate inputs.
	if err := ud.validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Store userData in context to be used by next handler
	ctx.Set("userData", ud)
	ctx.Next()
}
