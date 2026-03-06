package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/gowa/saas-openclaw/backend/internal/api/health"
	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/config"
	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/database"
	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/platform"
	customLogger "github.com/gowa/saas-openclaw/backend/pkg/logger"
	"github.com/gowa/saas-openclaw/backend/pkg/middleware"
	"github.com/gowa/saas-openclaw/backend/pkg/jwt"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := customLogger.New(cfg.Log.Level)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	// Initialize database connection
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer func() {
		if err := database.Close(db); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	logger.Info("Database connection established")

	// Initialize JWT validator
	var jwtValidator *jwt.Validator
	if cfg.JWT.Secret != "" {
		jwtValidator = jwt.NewValidator(cfg.JWT.Secret)
		logger.Info("JWT validator initialized")
	} else {
		logger.Warn("JWT_SECRET not configured, authentication middleware will not work")
	}

	// Initialize platform client
	var platformClient *platform.Client
	if cfg.Platform.BaseURL != "" {
		platformClient = platform.NewClient(cfg.Platform.BaseURL,
			platform.WithTimeout(cfg.Platform.Timeout),
		)
		logger.Info("Platform client initialized",
			zap.String("baseURL", cfg.Platform.BaseURL),
		)
	} else {
		logger.Warn("PLATFORM_BASE_URL not configured, platform integration will not work")
	}

	// Initialize Gin router
	router := gin.New()
	router.Use(middleware.Logger(logger), middleware.Recovery(logger))

	// Health check endpoints
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Database health check endpoint
	healthHandler := health.NewHandler(db, &cfg.Database, logger)
	router.GET("/health/database", healthHandler.Database)

	// API v1 routes
	v1 := router.Group("/v1")
	{
		// Public routes
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	// Protected routes (require authentication)
	if jwtValidator != nil && platformClient != nil {
		protected := router.Group("/v1")
		protected.Use(middleware.PlatformAuth(jwtValidator, platformClient))
		{
			// Get current user info
			protected.GET("/me", func(c *gin.Context) {
				user := middleware.GetCurrentUser(c)
				if user == nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"data": nil,
						"error": gin.H{
							"code":    "UNAUTHORIZED",
							"message": "User not authenticated",
						},
						"meta": nil,
					})
					return
				}
				c.JSON(http.StatusOK, gin.H{
					"data": gin.H{
						"id":       user.ID,
						"name":     user.Name,
						"email":    user.Email,
						"tenantId": user.TenantID,
						"role":     user.Role,
					},
					"error": nil,
					"meta":  nil,
				})
			})
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("Starting server",
			zap.String("port", cfg.Server.Port),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}
