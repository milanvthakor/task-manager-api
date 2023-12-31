package auth

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/milanvthakor/task-manager-api/internal/models"
	"github.com/milanvthakor/task-manager-api/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

// userData holds the authentication details of the user.
type userData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterHandler handles user registration.
func RegisterHandler(ctx *gin.Context, app *config.Application) {
	userData := ctx.MustGet("userData").(userData)

	// Check if the user with the email already exists in the database.
	user, err := app.UserRepository.GetUserByEmail(userData.Email)
	if err != nil {
		log.Printf("Warning: Failed to get user details from the database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify the email"})
		return
	}
	if user != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Hash password.
	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Warning: Failed to generate hash from the password: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process the password"})
		return
	}

	// Store user details in the database.
	user = &models.User{
		Email:    userData.Email,
		Password: string(hash),
	}
	if err := app.UserRepository.CreateUser(user); err != nil {
		log.Printf("Warning: Failed to register user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginHandler handles user login.
func LoginHandler(ctx *gin.Context, app *config.Application) {
	userData := ctx.MustGet("userData").(userData)

	// Check if the user with the email already exists in the database.
	user, err := app.UserRepository.GetUserByEmail(userData.Email)
	if err != nil {
		log.Printf("Warning: Failed to get user details from the database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify the email"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Email doesn't exists"})
		return
	}

	// Compare hashed password with the input password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userData.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong Password"})
		return
	}

	// Generate a JWT token.
	token, err := generateToken(app.Config.SecretKey, user.Email, user.ID)
	if err != nil {
		log.Printf("Warning: Failed to generate token: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
