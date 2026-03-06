// Package main provides CLI tools for admin user management
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/term"

	"github.com/gowa/saas-openclaw/backend/internal/domain/user"
	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/config"
	"github.com/gowa/saas-openclaw/backend/internal/infrastructure/database"
	"github.com/gowa/saas-openclaw/backend/internal/repository"
	customLogger "github.com/gowa/saas-openclaw/backend/pkg/logger"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 128
	maxUsernameLength = 50
	maxNameLength     = 100
	maxEmailLength    = 255
)

// emailRegex validates email format
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create-admin":
		if err := createAdmin(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Admin CLI - Platform Administration Tools")
	fmt.Println()
	fmt.Println("Usage: admin <command>")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  create-admin  Create a new admin user")
	fmt.Println("  help          Show this help message")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  DB_HOST       Database host (default: localhost)")
	fmt.Println("  DB_PORT       Database port (default: 5432)")
	fmt.Println("  DB_USER       Database user (default: postgres)")
	fmt.Println("  DB_PASSWORD   Database password")
	fmt.Println("  DB_NAME       Database name (default: saas_openclaw)")
	fmt.Println("  DB_SSLMODE    SSL mode (default: disable)")
	fmt.Println("  LOG_LEVEL     Log level (default: info)")
}

func createAdmin() error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger
	logger, err := customLogger.New(cfg.Log.Level)
	if err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}
	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err := database.Close(db); err != nil {
			logger.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	logger.Info("Connected to database")

	// Get input from user
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Println("=== Create Admin User ===")
	fmt.Println()

	// Username
	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}
	username = strings.TrimSpace(username)

	if username == "" {
		return fmt.Errorf("username is required")
	}
	if len(username) > maxUsernameLength {
		return fmt.Errorf("username must be at most %d characters", maxUsernameLength)
	}

	// Password
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	fmt.Println() // New line after password input

	if len(passwordBytes) == 0 {
		return fmt.Errorf("password is required")
	}
	if len(passwordBytes) < minPasswordLength {
		return fmt.Errorf("password must be at least %d characters", minPasswordLength)
	}
	if len(passwordBytes) > maxPasswordLength {
		return fmt.Errorf("password must be at most %d characters", maxPasswordLength)
	}

	// Confirm password
	fmt.Print("Confirm password: ")
	confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return fmt.Errorf("failed to read password confirmation: %w", err)
	}
	fmt.Println()

	if string(passwordBytes) != string(confirmBytes) {
		return fmt.Errorf("passwords do not match")
	}

	// Name
	fmt.Print("Enter name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read name: %w", err)
	}
	name = strings.TrimSpace(name)

	if name == "" {
		return fmt.Errorf("name is required")
	}
	if len(name) > maxNameLength {
		return fmt.Errorf("name must be at most %d characters", maxNameLength)
	}

	// Email
	fmt.Print("Enter email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read email: %w", err)
	}
	email = strings.TrimSpace(email)

	if email == "" {
		return fmt.Errorf("email is required")
	}
	if len(email) > maxEmailLength {
		return fmt.Errorf("email must be at most %d characters", maxEmailLength)
	}
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("invalid email format")
	}

	// Role
	fmt.Print("Enter role [admin/super_admin] (default: admin): ")
	roleInput, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read role: %w", err)
	}
	roleInput = strings.TrimSpace(roleInput)

	var role user.AdminRole
	switch strings.ToLower(roleInput) {
	case "super_admin":
		role = user.AdminRoleSuperAdmin
	case "", "admin":
		role = user.AdminRoleAdmin
	default:
		return fmt.Errorf("invalid role: %s (valid: admin, super_admin)", roleInput)
	}

	// Create admin user
	admin := &user.AdminUser{
		ID:       uuid.New().String(),
		Username: username,
		Name:     name,
		Email:    email,
		Role:     role,
	}

	// Save to database
	repo := repository.NewAdminUserRepository(db)
	if err := repo.Create(admin, string(passwordBytes)); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	fmt.Println()
	fmt.Println("=== Admin User Created Successfully ===")
	fmt.Printf("ID:       %s\n", admin.ID)
	fmt.Printf("Username: %s\n", admin.Username)
	fmt.Printf("Name:     %s\n", admin.Name)
	fmt.Printf("Email:    %s\n", admin.Email)
	fmt.Printf("Role:     %s\n", admin.Role)
	fmt.Printf("Created:  %s\n", admin.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	logger.Info("Admin user created",
		zap.String("adminId", admin.ID),
		zap.String("username", admin.Username),
		zap.String("role", string(admin.Role)),
	)

	return nil
}
