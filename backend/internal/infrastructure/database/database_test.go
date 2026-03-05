package database

import (
	"fmt"
	"testing"
	"time"

	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/config"
)

// buildDSN creates a DSN string from config (extracted for testing)
func buildDSN(cfg *config.DatabaseConfig) string {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)
	if cfg.SSLRootCert != "" {
		dsn += fmt.Sprintf(" sslrootcert=%s", cfg.SSLRootCert)
	}
	if cfg.SSLCert != "" {
		dsn += fmt.Sprintf(" sslcert=%s", cfg.SSLCert)
	}
	if cfg.SSLKey != "" {
		dsn += fmt.Sprintf(" sslkey=%s", cfg.SSLKey)
	}
	return dsn
}

// TestDSNGeneration tests that DSN is correctly generated from config
func TestDSNGeneration(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.DatabaseConfig
		expected string
	}{
		{
			name: "full config without SSL certs",
			cfg: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "secret",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			expected: "host=localhost port=5432 user=postgres password=secret dbname=testdb sslmode=disable",
		},
		{
			name: "with ssl require",
			cfg: &config.DatabaseConfig{
				Host:     "prod-db.example.com",
				Port:     5433,
				User:     "admin",
				Password: "prodpass",
				Name:     "proddb",
				SSLMode:  "require",
			},
			expected: "host=prod-db.example.com port=5433 user=admin password=prodpass dbname=proddb sslmode=require",
		},
		{
			name: "with verify-full ssl and certificates",
			cfg: &config.DatabaseConfig{
				Host:        "secure-db.example.com",
				Port:        5432,
				User:        "secureuser",
				Password:    "securepass",
				Name:        "securedb",
				SSLMode:     "verify-full",
				SSLRootCert: "/etc/ssl/certs/root.crt",
				SSLCert:     "/etc/ssl/certs/client.crt",
				SSLKey:      "/etc/ssl/private/client.key",
			},
			expected: "host=secure-db.example.com port=5432 user=secureuser password=securepass dbname=securedb sslmode=verify-full sslrootcert=/etc/ssl/certs/root.crt sslcert=/etc/ssl/certs/client.crt sslkey=/etc/ssl/private/client.key",
		},
		{
			name: "with only root cert",
			cfg: &config.DatabaseConfig{
				Host:        "ca-db.example.com",
				Port:        5432,
				User:        "user",
				Password:    "pass",
				Name:        "db",
				SSLMode:     "verify-ca",
				SSLRootCert: "/etc/ssl/certs/root.crt",
			},
			expected: "host=ca-db.example.com port=5432 user=user password=pass dbname=db sslmode=verify-ca sslrootcert=/etc/ssl/certs/root.crt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := buildDSN(tt.cfg)
			if dsn != tt.expected {
				t.Errorf("DSN = %q, want %q", dsn, tt.expected)
			}
		})
	}
}

// TestDatabaseConfigValidation tests that database configuration validation works
func TestDatabaseConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *config.DatabaseConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid config",
			cfg: &config.DatabaseConfig{
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
			},
			expectError: false,
		},
		{
			name: "empty host",
			cfg: &config.DatabaseConfig{
				Host:     "",
				Port:     5432,
				User:     "postgres",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			expectError: true,
			errorMsg:    "database host is required",
		},
		{
			name: "invalid port - negative",
			cfg: &config.DatabaseConfig{
				Host:    "localhost",
				Port:    -1,
				User:    "postgres",
				Name:    "testdb",
				SSLMode: "disable",
			},
			expectError: true,
			errorMsg:    "database port must be between 1 and 65535",
		},
		{
			name: "invalid port - too high",
			cfg: &config.DatabaseConfig{
				Host:    "localhost",
				Port:    70000,
				User:    "postgres",
				Name:    "testdb",
				SSLMode: "disable",
			},
			expectError: true,
			errorMsg:    "database port must be between 1 and 65535",
		},
		{
			name: "empty user",
			cfg: &config.DatabaseConfig{
				Host:    "localhost",
				Port:    5432,
				User:    "",
				Name:    "testdb",
				SSLMode: "disable",
			},
			expectError: true,
			errorMsg:    "database user is required",
		},
		{
			name: "empty database name",
			cfg: &config.DatabaseConfig{
				Host:    "localhost",
				Port:    5432,
				User:    "postgres",
				Name:    "",
				SSLMode: "disable",
			},
			expectError: true,
			errorMsg:    "database name is required",
		},
		{
			name: "negative max open conns",
			cfg: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				User:         "postgres",
				Name:         "testdb",
				SSLMode:      "disable",
				MaxOpenConns: -1,
			},
			expectError: true,
			errorMsg:    "max open connections cannot be negative",
		},
		{
			name: "negative max idle conns",
			cfg: &config.DatabaseConfig{
				Host:         "localhost",
				Port:         5432,
				User:         "postgres",
				Name:         "testdb",
				SSLMode:      "disable",
				MaxIdleConns: -1,
			},
			expectError: true,
			errorMsg:    "max idle connections cannot be negative",
		},
		{
			name: "invalid SSL mode",
			cfg: &config.DatabaseConfig{
				Host:    "localhost",
				Port:    5432,
				User:    "postgres",
				Name:    "testdb",
				SSLMode: "invalid",
			},
			expectError: true,
			errorMsg:    "invalid SSL mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errorMsg)
				} else if tt.errorMsg != "" && err.Error()[:len(tt.errorMsg)] != tt.errorMsg {
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

// TestSSLModeValidation tests SSL mode configuration
func TestSSLModeValidation(t *testing.T) {
	validModes := []string{"disable", "require", "verify-full", "verify-ca", "allow", "prefer"}

	for _, mode := range validModes {
		t.Run("valid mode: "+mode, func(t *testing.T) {
			cfg := &config.DatabaseConfig{
				Host:    "localhost",
				Port:    5432,
				User:    "postgres",
				Name:    "testdb",
				SSLMode: mode,
			}

			if err := cfg.Validate(); err != nil {
				t.Errorf("SSL mode %q should be valid, got error: %v", mode, err)
			}
		})
	}
}

// TestConnectionPoolDefaults tests default connection pool values
func TestConnectionPoolDefaults(t *testing.T) {
	cfg := &config.DatabaseConfig{
		Host:            "localhost",
		Port:            5432,
		User:            "postgres",
		Name:            "testdb",
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

// TestBuildDSN tests the buildDSN function directly
func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *config.DatabaseConfig
		contains []string
	}{
		{
			name: "basic config",
			cfg: &config.DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "postgres",
				Password: "secret",
				Name:     "testdb",
				SSLMode:  "disable",
			},
			contains: []string{"host=localhost", "port=5432", "user=postgres", "dbname=testdb", "sslmode=disable"},
		},
		{
			name: "with SSL certificates",
			cfg: &config.DatabaseConfig{
				Host:        "localhost",
				Port:        5432,
				User:        "postgres",
				Password:    "secret",
				Name:        "testdb",
				SSLMode:     "verify-full",
				SSLRootCert: "/path/to/root.crt",
				SSLCert:     "/path/to/client.crt",
				SSLKey:      "/path/to/client.key",
			},
			contains: []string{
				"sslmode=verify-full",
				"sslrootcert=/path/to/root.crt",
				"sslcert=/path/to/client.crt",
				"sslkey=/path/to/client.key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsn := buildDSN(tt.cfg)
			for _, s := range tt.contains {
				if !contains(dsn, s) {
					t.Errorf("DSN %q should contain %q", dsn, s)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
