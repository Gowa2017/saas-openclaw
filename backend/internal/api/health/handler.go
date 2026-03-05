package health

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"

	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/config"
)

// Handler handles health check requests
type Handler struct {
	db  *sqlx.DB
	cfg *config.DatabaseConfig
}

// NewHandler creates a new health handler
func NewHandler(db *sqlx.DB, cfg *config.DatabaseConfig) *Handler {
	return &Handler{db: db, cfg: cfg}
}

// PoolStats represents database connection pool statistics
type PoolStats struct {
	OpenConnections int `json:"open_connections"`
	InUse           int `json:"in_use"`
	Idle            int `json:"idle"`
	MaxOpen         int `json:"max_open"`
	MaxIdle         int `json:"max_idle"`
}

// DatabaseHealth represents database health status
type DatabaseHealth struct {
	Connected bool      `json:"connected"`
	Version   string    `json:"version"`
	PoolStats PoolStats `json:"pool_stats"`
}

// HealthData represents the data section of health response
type HealthData struct {
	Status   string         `json:"status"`
	Database DatabaseHealth `json:"database"`
}

// Meta represents metadata in the response
type Meta struct {
	Timestamp string `json:"timestamp"`
}

// Response represents the standard API response format
type Response struct {
	Data  HealthData `json:"data"`
	Error *string    `json:"error"`
	Meta  Meta       `json:"meta"`
}

// Database handles GET /health/database
func (h *Handler) Database(c *gin.Context) {
	var (
		connected bool
		version   string
		stats     sql.DBStats
		err       error
	)

	// Check database connection and get version
	if h.db != nil {
		err = h.db.Ping()
		if err == nil {
			connected = true
			// Query PostgreSQL version
			row := h.db.QueryRow("SELECT version()")
			if err := row.Scan(&version); err != nil {
				version = "unknown"
			}
			stats = h.db.Stats()
		}
	}

	// Determine status
	status := "healthy"
	var errorMsg *string
	if !connected {
		status = "unhealthy"
		msg := "database connection failed"
		if err != nil {
			msg = err.Error()
		}
		errorMsg = &msg
	}

	// Get MaxIdleConns from config, default to 0 if config is nil
	maxIdleConns := 0
	if h.cfg != nil {
		maxIdleConns = h.cfg.MaxIdleConns
	}

	// Build response
	response := Response{
		Data: HealthData{
			Status: status,
			Database: DatabaseHealth{
				Connected: connected,
				Version:   version,
				PoolStats: PoolStats{
					OpenConnections: stats.OpenConnections,
					InUse:           stats.InUse,
					Idle:            stats.Idle,
					MaxOpen:         stats.MaxOpenConnections,
					MaxIdle:         maxIdleConns,
				},
			},
		},
		Error: errorMsg,
		Meta: Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Set appropriate HTTP status code
	httpStatus := http.StatusOK
	if !connected {
		httpStatus = http.StatusServiceUnavailable
	}

	c.JSON(httpStatus, response)
}
