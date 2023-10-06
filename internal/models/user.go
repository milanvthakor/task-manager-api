package models

import "database/sql"

// User represents a user in the application.
type User struct {
	Id       uint   `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserRepository provides an interface for user-related database operations.
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user into the database.
func (r *UserRepository) CreateUser(user *User) error {
	_, err := r.db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	return err
}

// GetUserByEmail retrieves a user by email from the database.
func (r *UserRepository) GetUserByEmail(email string) (*User, error) {
	var user User
	err := r.db.QueryRow("SELECT * FROM users WHERE email = $1", email).Scan(&user.Id, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return nil, nil // Return nil when no records are found
	}
	if err != nil {
		return nil, err
	}

	return &user, nil
}
