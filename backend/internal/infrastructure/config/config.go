package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// ServerConfig holds HTTP server configuration
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	Name            string
	SSLMode         string
	SSLRootCert     string // Path to SSL root certificate (for verify-full)
	SSLCert         string // Path to SSL client certificate
	SSLKey          string // Path to SSL client key
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// Load reads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set defaults
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", 5432)
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_NAME", "saas_openclaw")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("DB_MAX_OPEN_CONNS", 100)
	viper.SetDefault("DB_MAX_IDLE_CONNS", 10)
	viper.SetDefault("DB_CONN_MAX_LIFETIME", "30m")
	viper.SetDefault("DB_CONN_MAX_IDLE_TIME", "10m")
	viper.SetDefault("DB_SSL_ROOT_CERT", "")
	viper.SetDefault("DB_SSL_CERT", "")
	viper.SetDefault("DB_SSL_KEY", "")
	viper.SetDefault("LOG_LEVEL", "info")

	// Read from environment variables
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, using defaults and env vars
	}

	cfg := &Config{
		Server: ServerConfig{
			Port:         viper.GetString("SERVER_PORT"),
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		Database: DatabaseConfig{
			Host:            viper.GetString("DB_HOST"),
			Port:            viper.GetInt("DB_PORT"),
			User:            viper.GetString("DB_USER"),
			Password:        viper.GetString("DB_PASSWORD"),
			Name:            viper.GetString("DB_NAME"),
			SSLMode:         viper.GetString("DB_SSLMODE"),
			SSLRootCert:     viper.GetString("DB_SSL_ROOT_CERT"),
			SSLCert:         viper.GetString("DB_SSL_CERT"),
			SSLKey:          viper.GetString("DB_SSL_KEY"),
			MaxOpenConns:    viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DB_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DB_CONN_MAX_LIFETIME"),
			ConnMaxIdleTime: viper.GetDuration("DB_CONN_MAX_IDLE_TIME"),
		},
		Log: LogConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	return c.Database.Validate()
}

// Validate validates the database configuration
func (d *DatabaseConfig) Validate() error {
	if d.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if d.Port <= 0 || d.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535, got %d", d.Port)
	}
	if d.User == "" {
		return fmt.Errorf("database user is required")
	}
	if d.Name == "" {
		return fmt.Errorf("database name is required")
	}
	if d.MaxOpenConns < 0 {
		return fmt.Errorf("max open connections cannot be negative")
	}
	if d.MaxIdleConns < 0 {
		return fmt.Errorf("max idle connections cannot be negative")
	}

	validSSLModes := map[string]bool{
		"disable":     true,
		"allow":       true,
		"prefer":      true,
		"require":     true,
		"verify-ca":   true,
		"verify-full": true,
	}
	if d.SSLMode != "" && !validSSLModes[d.SSLMode] {
		return fmt.Errorf("invalid SSL mode: %s (valid: disable, allow, prefer, require, verify-ca, verify-full)", d.SSLMode)
	}

	return nil
}
