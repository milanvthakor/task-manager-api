package config

import "github.com/milanvthakor/task-manager-api/internal/models"

// Application holds application-wide dependencies.
type Application struct {
	UserRepository models.UserRepository
}
