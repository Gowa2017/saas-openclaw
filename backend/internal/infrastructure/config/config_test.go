package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("SERVER_PORT", "9090")
	os.Setenv("DB_HOST", "testhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSLMODE", "require")
	os.Setenv("LOG_LEVEL", "debug")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("LOG_LEVEL")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("Server.Port = %v, want 9090", cfg.Server.Port)
	}

	if cfg.Database.Host != "testhost" {
		t.Errorf("Database.Host = %v, want testhost", cfg.Database.Host)
	}

	if cfg.Database.Port != 5433 {
		t.Errorf("Database.Port = %v, want 5433", cfg.Database.Port)
	}

	if cfg.Database.SSLMode != "require" {
		t.Errorf("Database.SSLMode = %v, want require", cfg.Database.SSLMode)
	}

	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %v, want debug", cfg.Log.Level)
	}
}

func TestLoadDefaults(t *testing.T) {
	// Clear all relevant env vars
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("DB_SSLMODE")
	os.Unsetenv("LOG_LEVEL")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("Server.Port default = %v, want 8080", cfg.Server.Port)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("Database.Host default = %v, want localhost", cfg.Database.Host)
	}

	if cfg.Database.SSLMode != "disable" {
		t.Errorf("Database.SSLMode default = %v, want disable", cfg.Database.SSLMode)
	}

	if cfg.Database.MaxOpenConns != 25 {
		t.Errorf("Database.MaxOpenConns default = %v, want 25", cfg.Database.MaxOpenConns)
	}

	if cfg.Database.MaxIdleConns != 5 {
		t.Errorf("Database.MaxIdleConns default = %v, want 5", cfg.Database.MaxIdleConns)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level default = %v, want info", cfg.Log.Level)
	}
}
