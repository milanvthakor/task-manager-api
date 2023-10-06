package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/milanvthakor/task-manager-api/internal/validator"
	"github.com/milanvthakor/task-manager-api/pkg/config"
)

// ValidateInputMiddleware validates the authentication inputs.
func ValidateInputMiddleware(ctx *gin.Context, app *config.Application) {
	var ud userData
	if err := ctx.ShouldBindJSON(&ud); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid inputs"})
		return
	}

	// Validate inputs
	if !validator.IsValidEmail(ud.Email) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
		return
	}
	if !validator.IsValidPassword(ud.Password) {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid password. The length must be between 8 and 12 characters"})
		return
	}

	// Store userData in context to be used by next handler
	ctx.Set("userData", ud)
	ctx.Next()
}

// AuthenticateMiddleware authenticates the incoming request by validating JWT token.
func AuthenticateMiddleware(ctx *gin.Context, app *config.Application) {
	// Get the authorization header
	authHeader := ctx.GetHeader("Authorization")
	// Extract the JWT token from the authorization header
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
		return
	}

	// Validate token
	claims, err := validateToken(app.Config.SecretKey, token)
	if err != nil {
		log.Printf("Warning: Failed to validate JWT token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// Get the userID from the token and store in the context for later use
	userID := claims["userID"].(float64)
	ctx.Set("userID", uint(userID))

	ctx.Next()
}
