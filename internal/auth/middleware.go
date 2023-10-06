package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/milanvthakor/task-manager-api/pkg/config"
)

// Middleware for authentication input validation
func ValidateInput(ctx *gin.Context, app *config.Application) {
	var ud userData
	if err := ctx.ShouldBindJSON(&ud); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid inputs"})
		return
	}

	// Validate inputs.
	if err := ud.validate(); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Store userData in context to be used by next handler
	ctx.Set("userData", ud)
	ctx.Next()
}

// Authenticate authenticates the incoming request by validating JWT token.
func Authenticate(ctx *gin.Context, app *config.Application) {
	// Get the authorization header
	authHeader := ctx.GetHeader("Authorization")
	// Extract the JWT token from the authorization header
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
		return
	}

	// Validate token
	claims, err := validateToken(app.Config.SecretKey, token)
	if err != nil {
		log.Printf("Warning: Failed to validate JWT token: %v", err)
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	// Get the userID from the token and store in the context for later use
	userID := claims["userID"].(float64)
	ctx.Set("userID", uint(userID))

	ctx.Next()
}

// validateToken validates a JWT token.
func validateToken(secretKey, tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}
