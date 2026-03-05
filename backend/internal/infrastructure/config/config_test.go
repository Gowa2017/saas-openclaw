package config

import (
	"os"
	"testing"
	"time"
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
	os.Setenv("DB_MAX_OPEN_CONNS", "50")
	os.Setenv("DB_MAX_IDLE_CONNS", "15")
	os.Setenv("DB_CONN_MAX_LIFETIME", "1h")
	os.Setenv("DB_CONN_MAX_IDLE_TIME", "30m")
	os.Setenv("DB_SSL_ROOT_CERT", "/path/to/root.crt")
	os.Setenv("DB_SSL_CERT", "/path/to/client.crt")
	os.Setenv("DB_SSL_KEY", "/path/to/client.key")
	os.Setenv("LOG_LEVEL", "debug")

	defer func() {
		os.Unsetenv("SERVER_PORT")
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASSWORD")
		os.Unsetenv("DB_NAME")
		os.Unsetenv("DB_SSLMODE")
		os.Unsetenv("DB_MAX_OPEN_CONNS")
		os.Unsetenv("DB_MAX_IDLE_CONNS")
		os.Unsetenv("DB_CONN_MAX_LIFETIME")
		os.Unsetenv("DB_CONN_MAX_IDLE_TIME")
		os.Unsetenv("DB_SSL_ROOT_CERT")
		os.Unsetenv("DB_SSL_CERT")
		os.Unsetenv("DB_SSL_KEY")
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

	if cfg.Database.MaxOpenConns != 50 {
		t.Errorf("Database.MaxOpenConns = %v, want 50", cfg.Database.MaxOpenConns)
	}

	if cfg.Database.MaxIdleConns != 15 {
		t.Errorf("Database.MaxIdleConns = %v, want 15", cfg.Database.MaxIdleConns)
	}

	if cfg.Database.ConnMaxLifetime != time.Hour {
		t.Errorf("Database.ConnMaxLifetime = %v, want 1h", cfg.Database.ConnMaxLifetime)
	}

	if cfg.Database.ConnMaxIdleTime != 30*time.Minute {
		t.Errorf("Database.ConnMaxIdleTime = %v, want 30m", cfg.Database.ConnMaxIdleTime)
	}

	if cfg.Database.SSLRootCert != "/path/to/root.crt" {
		t.Errorf("Database.SSLRootCert = %v, want /path/to/root.crt", cfg.Database.SSLRootCert)
	}

	if cfg.Database.SSLCert != "/path/to/client.crt" {
		t.Errorf("Database.SSLCert = %v, want /path/to/client.crt", cfg.Database.SSLCert)
	}

	if cfg.Database.SSLKey != "/path/to/client.key" {
		t.Errorf("Database.SSLKey = %v, want /path/to/client.key", cfg.Database.SSLKey)
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
	os.Unsetenv("DB_MAX_OPEN_CONNS")
	os.Unsetenv("DB_MAX_IDLE_CONNS")
	os.Unsetenv("DB_CONN_MAX_LIFETIME")
	os.Unsetenv("DB_CONN_MAX_IDLE_TIME")
	os.Unsetenv("DB_SSL_ROOT_CERT")
	os.Unsetenv("DB_SSL_CERT")
	os.Unsetenv("DB_SSL_KEY")
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

	if cfg.Database.MaxOpenConns != 100 {
		t.Errorf("Database.MaxOpenConns default = %v, want 100", cfg.Database.MaxOpenConns)
	}

	if cfg.Database.MaxIdleConns != 10 {
		t.Errorf("Database.MaxIdleConns default = %v, want 10", cfg.Database.MaxIdleConns)
	}

	if cfg.Database.ConnMaxLifetime != 30*time.Minute {
		t.Errorf("Database.ConnMaxLifetime default = %v, want 30m", cfg.Database.ConnMaxLifetime)
	}

	if cfg.Database.ConnMaxIdleTime != 10*time.Minute {
		t.Errorf("Database.ConnMaxIdleTime default = %v, want 10m", cfg.Database.ConnMaxIdleTime)
	}

	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level default = %v, want info", cfg.Log.Level)
	}

	// SSL cert defaults should be empty
	if cfg.Database.SSLRootCert != "" {
		t.Errorf("Database.SSLRootCert default = %v, want empty", cfg.Database.SSLRootCert)
	}
	if cfg.Database.SSLCert != "" {
		t.Errorf("Database.SSLCert default = %v, want empty", cfg.Database.SSLCert)
	}
	if cfg.Database.SSLKey != "" {
		t.Errorf("Database.SSLKey default = %v, want empty", cfg.Database.SSLKey)
	}
}

func TestDatabaseConfigValidate(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *DatabaseConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			cfg: &DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				User:         "postgres",
				Name:         "testdb",
				SSLMode:      "disable",
				MaxOpenConns: 100,
				MaxIdleConns: 10,
			},
			expectError: false,
		},
		{
			name: "empty host",
			cfg: &DatabaseConfig{
				Host: "",
				Port: 5432,
				User: "postgres",
				Name: "testdb",
			},
			expectError: true,
			errorMsg:    "database host is required",
		},
		{
			name: "invalid port - negative",
			cfg: &DatabaseConfig{
				Host: "localhost",
				Port: -1,
				User: "postgres",
				Name: "testdb",
			},
			expectError: true,
			errorMsg:    "database port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			cfg: &DatabaseConfig{
				Host: "localhost",
				Port: 70000,
				User: "postgres",
				Name: "testdb",
			},
			expectError: true,
			errorMsg:    "database port must be between 1 and 65535",
		},
		{
			name: "empty user",
			cfg: &DatabaseConfig{
				Host: "localhost",
				Port: 5432,
				User: "",
				Name: "testdb",
			},
			expectError: true,
			errorMsg:    "database user is required",
		},
		{
			name: "empty name",
			cfg: &DatabaseConfig{
				Host: "localhost",
				Port: 5432,
				User: "postgres",
				Name: "",
			},
			expectError: true,
			errorMsg:    "database name is required",
		},
		{
			name: "negative max open conns",
			cfg: &DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				User:         "postgres",
				Name:         "testdb",
				MaxOpenConns: -1,
			},
			expectError: true,
			errorMsg:    "max open connections cannot be negative",
		},
		{
			name: "negative max idle conns",
			cfg: &DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				User:         "postgres",
				Name:         "testdb",
				MaxIdleConns: -1,
			},
			expectError: true,
			errorMsg:    "max idle connections cannot be negative",
		},
		{
			name: "invalid SSL mode",
			cfg: &DatabaseConfig{
				Host:    "localhost",
				Port:    5432,
				User:    "postgres",
				Name:    "testdb",
				SSLMode: "invalid",
			},
			expectError: true,
			errorMsg:    "invalid SSL mode",
		},
		{
			name: "valid SSL modes",
			cfg: &DatabaseConfig{
				Host:    "localhost",
				Port:    5432,
				User:    "postgres",
				Name:    "testdb",
				SSLMode: "verify-full",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if tt.errorMsg != "" && len(err.Error()) >= len(tt.errorMsg) && err.Error()[:len(tt.errorMsg)] != tt.errorMsg {
					t.Errorf("error = %q, want to contain %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestConfigValidate(t *testing.T) {
	// Test that Config.Validate delegates to DatabaseConfig.Validate
	cfg := &Config{
		Database: DatabaseConfig{
			Host: "", // Invalid
			Port: 5432,
			User: "postgres",
			Name: "testdb",
		},
	}

	err := cfg.Validate()
	if err == nil {
		t.Error("expected error for empty host")
	}

	// Valid config should not error
	cfg2 := &Config{
		Database: DatabaseConfig{
			Host: "localhost",
			Port: 5432,
			User: "postgres",
			Name: "testdb",
		},
	}

	if err := cfg2.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
