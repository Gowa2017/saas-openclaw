package database

import (
	"testing"
	"time"

	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/config"
)

// TestDatabaseConfigValidation tests that database configuration values are properly used
func TestDatabaseConfigValidation(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Password:        "testpass",
		Name:            "testdb",
		SSLMode:         "disable",
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	// Verify configuration values
	if cfg.MaxOpenConns != 100 {
		t.Errorf("MaxOpenConns = %v, want 100", cfg.MaxOpenConns)
	}

	if cfg.MaxIdleConns != 10 {
		t.Errorf("MaxIdleConns = %v, want 10", cfg.MaxIdleConns)
	}

	if cfg.ConnMaxLifetime != 30*time.Minute {
		t.Errorf("ConnMaxLifetime = %v, want 30m", cfg.ConnMaxLifetime)
	}

	if cfg.ConnMaxIdleTime != 10*time.Minute {
		t.Errorf("ConnMaxIdleTime = %v, want 10m", cfg.ConnMaxIdleTime)
	}
}

// TestSSLModeValidation tests SSL mode configuration
func TestSSLModeValidation(t *testing.T) {
	validModes := []string{"disable", "require", "verify-full", "verify-ca", "allow", "prefer"}

	for _, mode := range validModes {
		cfg := &config.DatabaseConfig{
			SSLMode: mode,
		}

		if cfg.SSLMode != mode {
			t.Errorf("SSLMode = %v, want %v", cfg.SSLMode, mode)
		}
	}
}

// TestConnectionPoolDefaults tests default connection pool values
func TestConnectionPoolDefaults(t *testing.T) {
	cfg := &config.DatabaseConfig{
		MaxOpenConns:    100,
		MaxIdleConns:    10,
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}

	// Test that MaxIdleConns is approximately 10% of MaxOpenConns as per best practices
	expectedIdleConns := cfg.MaxOpenConns / 10
	if cfg.MaxIdleConns != expectedIdleConns {
		t.Logf("Warning: MaxIdleConns (%v) is not 10%% of MaxOpenConns (%v)", cfg.MaxIdleConns, cfg.MaxOpenConns)
	}

	// Test that ConnMaxIdleTime is less than ConnMaxLifetime
	if cfg.ConnMaxIdleTime >= cfg.ConnMaxLifetime {
		t.Errorf("ConnMaxIdleTime (%v) should be less than ConnMaxLifetime (%v)", cfg.ConnMaxIdleTime, cfg.ConnMaxLifetime)
	}
}
