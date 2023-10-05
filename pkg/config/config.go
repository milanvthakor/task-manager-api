package config

import "os"

// Config holds the application configuration.
type Config struct {
	DatabaseDSN string
	ServerPort  string
}

// New creates a new Config instance with the default values.
func New() *Config {
	return &Config{
		DatabaseDSN: getEnv("DatabaseDSN", "postgres://username:password@localhost:5432/task-manager-database"),
		ServerPort:  getEnv("ServerPort", "8080"),
	}
}

// getEnv retrieves an environment variable with a default value if not set.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return value
}
