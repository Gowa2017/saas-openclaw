package database

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/config"
)

// Connect establishes a connection to the PostgreSQL database
func Connect(cfg *config.DatabaseConfig) (*sqlx.DB, error) {
	// Validate configuration
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid database configuration: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	// Add SSL certificate parameters if configured
	if cfg.SSLRootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", cfg.SSLRootCert)
	}
	if cfg.SSLCert != "" {
		dsn += fmt.Sprintf(" sslcert=%s", cfg.SSLCert)
	}
	if cfg.SSLKey != "" {
		dsn += fmt.Sprintf(" sslkey=%s", cfg.SSLKey)
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool from config
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

// Ping checks if the database connection is alive
func Ping(db *sqlx.DB) error {
	return db.Ping()
}

// Stats returns database connection pool statistics
func Stats(db *sqlx.DB) sql.DBStats {
	return db.Stats()
}

// Close closes the database connection
func Close(db *sqlx.DB) error {
	return db.Close()
}
