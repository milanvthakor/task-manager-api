package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milanvthakor/task-manager-api/internal/models"
	"github.com/milanvthakor/task-manager-api/internal/validator"
	"github.com/milanvthakor/task-manager-api/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

// UserData holds the authentication details of the user.
type UserData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler handles user registration
func RegisterHandler(app *config.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var ud UserData
		if err := ctx.ShouldBindJSON(&ud); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid inputs"})
			return
		}

		// Validate inputs
		if !validator.IsValidEmail(ud.Email) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid email"})
			return
		}
		if !validator.IsValidPassword(ud.Password) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid password. The length must be between 8 and 12 characters"})
			return
		}

		// Check if the user with the email already exists in the database.
		user, err := app.UserRepository.GetUserByEmail(ud.Email)
		if err != nil {
			log.Printf("Warning: Failed to get user details from the database: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to verify the email"})
			return
		}
		if user != nil {
			ctx.JSON(http.StatusConflict, gin.H{"message": "Email already exists"})
			return
		}

		// Hash password
		hash, err := bcrypt.GenerateFromPassword([]byte(ud.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Warning: Failed to generate hash from the password: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to process the password"})
			return
		}

		// Store user details in the database.
		user = &models.User{
			Email:    ud.Email,
			Password: string(hash),
		}
		if err := app.UserRepository.CreateUser(user); err != nil {
			log.Printf("Waiting: Failed to register user: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to register user"})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
	}
}
