package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/config"
)

func init() {
	gin.SetMode(gin.TestMode)
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
			// Create handler
			handler := NewHandler(nil, cfg)

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
			// MaxOpen comes from stats.MaxOpenConnections which is 0 when no DB
			// So we only verify MaxIdle comes from config
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

	handler := NewHandler(nil, cfg)

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
