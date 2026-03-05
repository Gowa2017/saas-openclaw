package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// getTestLogger returns a nop logger for testing
func getTestLogger() *zap.Logger {
	return zap.NewNop()
}

// TestDatabaseHealthHandler tests the database health endpoint
func TestDatabaseHealthHandler(t *testing.T) {
	// Create test config
	cfg := &config.DatabaseConfig{
		Host:         "localhost",
		Port:         5432,
		User:         "postgres",
		Password:     "testpass",
		Name:         "testdb",
		SSLMode:      "disable",
		MaxOpenConns: 100,
		MaxIdleConns: 10,
	}

	tests := []struct {
		name          string
		expectStatus  int
		expectHealthy bool
		withDB        bool
	}{
		{
			name:          "without database connection",
			expectStatus:  http.StatusServiceUnavailable,
			expectHealthy: false,
			withDB:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create handler with nop logger
			handler := NewHandler(nil, cfg, getTestLogger())

			// Create test router
			router := gin.New()
			router.GET("/health/database", handler.Database)

			// Create test request
			req, _ := http.NewRequest("GET", "/health/database", nil)
			w := httptest.NewRecorder()

			// Serve request
			router.ServeHTTP(w, req)

			// Check status code
			if w.Code != tt.expectStatus {
				t.Errorf("status code = %v, want %v", w.Code, tt.expectStatus)
			}

			// Parse response
			var response Response
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to parse response: %v", err)
			}

			// Check health status
			if response.Data.Status == "healthy" && !tt.expectHealthy {
				t.Errorf("expected unhealthy status, got healthy")
			}

			// Check that timestamp is present
			if response.Meta.Timestamp == "" {
				t.Error("timestamp is empty")
			}

			// Without DB connection, the pool stats come from config (MaxIdle)
			if response.Data.Database.PoolStats.MaxIdle != cfg.MaxIdleConns {
				t.Errorf("MaxIdle = %v, want %v", response.Data.Database.PoolStats.MaxIdle, cfg.MaxIdleConns)
			}

			// Verify connected is false when no DB
			if response.Data.Database.Connected {
				t.Error("expected connected to be false when no database")
			}

			// Verify error message is present
			if response.Error == nil {
				t.Error("expected error message when database is not connected")
			}
		})
	}
}

// TestHealthResponseStructure tests that the response structure matches the API spec
func TestHealthResponseStructure(t *testing.T) {
	cfg := &config.DatabaseConfig{
		MaxOpenConns: 100,
		MaxIdleConns: 10,
	}

	handler := NewHandler(nil, cfg, getTestLogger())

	router := gin.New()
	router.GET("/health/database", handler.Database)

	req, _ := http.NewRequest("GET", "/health/database", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}

	// Verify required top-level fields
	requiredFields := []string{"data", "error", "meta"}
	for _, field := range requiredFields {
		if _, exists := response[field]; !exists {
			t.Errorf("missing required field: %s", field)
		}
	}

	// Verify data structure
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("data field is not a map")
	}

	// Verify data.status exists
	if _, exists := data["status"]; !exists {
		t.Error("missing data.status field")
	}

	// Verify data.database exists
	dbData, ok := data["database"].(map[string]interface{})
	if !ok {
		t.Fatal("data.database field is not a map")
	}

	// Verify database fields
	dbRequiredFields := []string{"connected", "version", "pool_stats"}
	for _, field := range dbRequiredFields {
		if _, exists := dbData[field]; !exists {
			t.Errorf("missing data.database.%s field", field)
		}
	}

	// Verify pool_stats structure
	poolStats, ok := dbData["pool_stats"].(map[string]interface{})
	if !ok {
		t.Fatal("data.database.pool_stats field is not a map")
	}

	poolRequiredFields := []string{"open_connections", "in_use", "idle", "max_open", "max_idle"}
	for _, field := range poolRequiredFields {
		if _, exists := poolStats[field]; !exists {
			t.Errorf("missing data.database.pool_stats.%s field", field)
		}
	}

	// Verify meta structure
	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("meta field is not a map")
	}

	if _, exists := meta["timestamp"]; !exists {
		t.Error("missing meta.timestamp field")
	}
}

// TestNewHandler tests the NewHandler constructor
func TestNewHandler(t *testing.T) {
	cfg := &config.DatabaseConfig{
		MaxOpenConns: 100,
		MaxIdleConns: 10,
	}

	logger := getTestLogger()
	handler := NewHandler(nil, cfg, logger)
	if handler == nil {
		t.Fatal("NewHandler returned nil")
	}

	if handler.cfg != cfg {
		t.Error("handler config not set correctly")
	}

	if handler.db != nil {
		t.Error("handler db should be nil when nil is passed")
	}

	if handler.logger != logger {
		t.Error("handler logger not set correctly")
	}
}

// TestHandlerWithNilConfig tests handler behavior with nil config
func TestHandlerWithNilConfig(t *testing.T) {
	handler := NewHandler(nil, nil, getTestLogger())
	if handler == nil {
		t.Fatal("NewHandler should not return nil even with nil config")
	}

	router := gin.New()
	router.GET("/health/database", handler.Database)

	req, _ := http.NewRequest("GET", "/health/database", nil)
	w := httptest.NewRecorder()

	// Should not panic
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusServiceUnavailable)
	}
}

// TestHandlerWithNilLogger tests that handler works with nil logger
func TestHandlerWithNilLogger(t *testing.T) {
	cfg := &config.DatabaseConfig{
		MaxOpenConns: 100,
		MaxIdleConns: 10,
	}

	handler := NewHandler(nil, cfg, nil)
	if handler == nil {
		t.Fatal("NewHandler should not return nil")
	}

	router := gin.New()
	router.GET("/health/database", handler.Database)

	req, _ := http.NewRequest("GET", "/health/database", nil)
	w := httptest.NewRecorder()

	// Should not panic even with nil logger
	router.ServeHTTP(w, req)

	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusServiceUnavailable)
	}
}

// TestPoolStatsStructure tests PoolStats struct
func TestPoolStatsStructure(t *testing.T) {
	stats := PoolStats{
		OpenConnections: 10,
		InUse:           5,
		Idle:            5,
		MaxOpen:         100,
		MaxIdle:         10,
	}

	if stats.OpenConnections != 10 {
		t.Errorf("OpenConnections = %v, want 10", stats.OpenConnections)
	}
	if stats.InUse != 5 {
		t.Errorf("InUse = %v, want 5", stats.InUse)
	}
	if stats.Idle != 5 {
		t.Errorf("Idle = %v, want 5", stats.Idle)
	}
}

// TestDatabaseHealthStructure tests DatabaseHealth struct
func TestDatabaseHealthStructure(t *testing.T) {
	health := DatabaseHealth{
		Connected: true,
		Version:   "PostgreSQL 15.0",
		PoolStats: PoolStats{
			OpenConnections: 10,
			InUse:           5,
			Idle:            5,
			MaxOpen:         100,
			MaxIdle:         10,
		},
	}

	if !health.Connected {
		t.Error("Connected should be true")
	}
	if health.Version != "PostgreSQL 15.0" {
		t.Errorf("Version = %v, want PostgreSQL 15.0", health.Version)
	}
}

// TestHealthDataStructure tests HealthData struct
func TestHealthDataStructure(t *testing.T) {
	data := HealthData{
		Status: "healthy",
		Database: DatabaseHealth{
			Connected: true,
			Version:   "PostgreSQL 15.0",
		},
	}

	if data.Status != "healthy" {
		t.Errorf("Status = %v, want healthy", data.Status)
	}
}

// TestMetaStructure tests Meta struct
func TestMetaStructure(t *testing.T) {
	meta := Meta{
		Timestamp: "2024-01-01T00:00:00Z",
	}

	if meta.Timestamp != "2024-01-01T00:00:00Z" {
		t.Errorf("Timestamp = %v, want 2024-01-01T00:00:00Z", meta.Timestamp)
	}
}

// TestResponseStructure tests Response struct
func TestResponseStructure(t *testing.T) {
	errMsg := "test error"
	response := Response{
		Data: HealthData{
			Status: "healthy",
		},
		Error: &errMsg,
		Meta: Meta{
			Timestamp: "2024-01-01T00:00:00Z",
		},
	}

	if response.Data.Status != "healthy" {
		t.Errorf("Data.Status = %v, want healthy", response.Data.Status)
	}
	if response.Error == nil || *response.Error != "test error" {
		t.Error("Error not set correctly")
	}
}

// TestResponseWithNilError tests Response with nil error
func TestResponseWithNilError(t *testing.T) {
	response := Response{
		Data: HealthData{
			Status: "healthy",
		},
		Error: nil,
		Meta: Meta{
			Timestamp: "2024-01-01T00:00:00Z",
		},
	}

	if response.Error != nil {
		t.Error("Error should be nil")
	}
}
