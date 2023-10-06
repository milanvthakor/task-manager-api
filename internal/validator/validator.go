package validator

import (
	"regexp"
	"strings"

	"github.com/milanvthakor/task-manager-api/internal/models"
)

// emailRegex is a regular expression for validating email
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// IsValidEmail checks if an email is valid.
func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// IsValidPassword checks if a password is valid.
func IsValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 12
}

// IsBlank checks if a string is empty.
func IsBlank(s string) bool {
	s = strings.TrimSpace(s)
	return len(s) == 0
}

// IsValidTaskStatus checks if a status is valid.
func IsValidTaskStatus(status models.TaskStatus) bool {
	switch status {
	case models.TaskStatusTodo, models.TaskStatusInProgress, models.TaskStatusDone:
		return true
	}

	return false
}
