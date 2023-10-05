package validator

import "regexp"

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
