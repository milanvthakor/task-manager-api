package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

// Init initializes the database connection.
func Init(databaseDSN string) (*sql.DB, error) {
	var err error

	// Use sync.Once to ensure initialization only happens once.
	once.Do(func() {
		// Open a new database connection
		db, err = sql.Open("postgres", databaseDSN)
		if err != nil {
			err = fmt.Errorf("failed to open database connection: %w", err)
			return
		}

		// Test the database connection.
		if err = db.Ping(); err != nil {
			err = fmt.Errorf("failed to ping database: %w", err)
			return
		}

		log.Println("Connected to the database")
	})

	return db, err
}

// GetDB returns the database instance.
func GetDB() *sql.DB {
	return db
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
