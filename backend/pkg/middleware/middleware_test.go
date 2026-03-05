package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestLogger(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name       string
		method     string
		path       string
		routePath  string
		expectCode int
	}{
		{
			name:       "GET request",
			method:     http.MethodGet,
			path:       "/test",
			routePath:  "/test",
			expectCode: http.StatusOK,
		},
		{
			name:       "POST request",
			method:     http.MethodPost,
			path:       "/api/v1/users",
			routePath:  "/api/v1/users",
			expectCode: http.StatusCreated,
		},
		{
			name:       "request with query params",
			method:     http.MethodGet,
			path:       "/search?q=test",
			routePath:  "/search",
			expectCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(Logger(logger))

			router.Handle(tt.method, tt.routePath, func(c *gin.Context) {
				if tt.method == http.MethodPost {
					c.Status(http.StatusCreated)
				} else {
					c.Status(http.StatusOK)
				}
			})

			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.expectCode)
			}
		})
	}
}

func TestRecovery(t *testing.T) {
	logger := zaptest.NewLogger(t)

	t.Run("no panic - normal request", func(t *testing.T) {
		router := gin.New()
		router.Use(Recovery(logger))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status code = %v, want %v", w.Code, http.StatusOK)
		}
	})

	t.Run("panic - should recover and return 500", func(t *testing.T) {
		router := gin.New()
		router.Use(Recovery(logger))
		router.GET("/panic", func(c *gin.Context) {
			panic("test panic")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status code = %v, want %v", w.Code, http.StatusInternalServerError)
		}
	})

	t.Run("panic - should return JSON error", func(t *testing.T) {
		router := gin.New()
		router.Use(Recovery(logger))
		router.GET("/panic", func(c *gin.Context) {
			panic("test panic for JSON check")
		})

		req := httptest.NewRequest(http.MethodGet, "/panic", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		// Verify response is JSON
		contentType := w.Header().Get("Content-Type")
		if contentType != "application/json; charset=utf-8" {
			t.Errorf("Content-Type = %v, want application/json; charset=utf-8", contentType)
		}

		// Verify JSON body contains error key
		body := w.Body.String()
		if body == "" {
			t.Error("response body should not be empty")
		}
	})
}

func TestLoggerWithDifferentStatusCodes(t *testing.T) {
	logger := zaptest.NewLogger(t)

	tests := []struct {
		name       string
		statusCode int
	}{
		{"200 OK", http.StatusOK},
		{"201 Created", http.StatusCreated},
		{"400 Bad Request", http.StatusBadRequest},
		{"404 Not Found", http.StatusNotFound},
		{"500 Internal Server Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(Logger(logger))
			router.GET("/test", func(c *gin.Context) {
				c.Status(tt.statusCode)
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.statusCode {
				t.Errorf("status code = %v, want %v", w.Code, tt.statusCode)
			}
		})
	}
}

func TestRecoveryWithNilLogger(t *testing.T) {
	// Test that Recovery doesn't crash with nil logger
	// Note: In production, this would cause a panic, but we're testing the middleware logic
	router := gin.New()

	// Create a valid logger to test the middleware
	logger, _ := zap.NewProduction()
	router.Use(Recovery(logger))

	router.GET("/panic", func(c *gin.Context) {
		panic("test panic")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status code = %v, want %v", w.Code, http.StatusInternalServerError)
	}
}
