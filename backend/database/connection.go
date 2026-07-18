package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"jagapilar-backend/config"
)

var DB *sql.DB

// Connect establishes a connection to the Supabase PostgreSQL database
func Connect(cfg *config.Config) error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("✅ Connected to Supabase PostgreSQL")
	return nil
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
