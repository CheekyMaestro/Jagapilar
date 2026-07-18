package config

import (
	"os"
)

// Config holds all application configuration
type Config struct {
	// Server
	ServerPort string

	// Database (Supabase PostgreSQL)
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Supabase
	SupabaseURL string
	SupabaseKey string
}

// Load reads configuration from environment variables with defaults
func Load() *Config {
	return &Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),

		// Supabase PostgreSQL connection
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "jagapilar"),
		DBSSLMode:  getEnv("DB_SSLMODE", "require"), // Supabase requires SSL

		// Supabase API (for future use)
		SupabaseURL: getEnv("SUPABASE_URL", ""),
		SupabaseKey: getEnv("SUPABASE_KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
