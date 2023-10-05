package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Init initializes the database connection.
func Init(databaseDSN string) (*sql.DB, error) {
	var err error

	// Open a new database connection
	db, err = sql.Open("postgres", databaseDSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the database connection.
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to the database")
	return db, nil
}

// Close closes the database connection.
func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
		} else {
			log.Print("Database connection closed")
		}
	}
}
