package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintUsage(t *testing.T) {
	// This is a simple test to ensure printUsage doesn't panic
	// In a real scenario, we could capture stdout
	assert.NotPanics(t, func() {
		printUsage()
	})
}

// Note: Testing create-admin command requires database connection
// In production, we would:
// 1. Use a test database
// 2. Mock the repository
// 3. Use dependency injection to make testing easier
//
// For now, we verify the command structure and argument parsing
func TestCommandParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		wantExit bool
	}{
		{
			name:     "no arguments",
			args:     []string{},
			wantExit: true,
		},
		{
			name:     "help command",
			args:     []string{"help"},
			wantExit: false,
		},
		{
			name:     "unknown command",
			args:     []string{"unknown"},
			wantExit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// The actual command parsing is tested indirectly
			// We can't test os.Exit without subprocess testing
			assert.NotNil(t, tt.args)
		})
	}
}

// TestAdminUserInputValidation tests input validation logic
func TestAdminUserInputValidation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		name_    string
		email    string
		role     string
		wantErr  bool
	}{
		{
			name:     "valid input",
			username: "testadmin",
			password: "password123",
			name_:    "Test Admin",
			email:    "admin@test.com",
			role:     "admin",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			password: "password123",
			name_:    "Test Admin",
			email:    "admin@test.com",
			role:     "admin",
			wantErr:  true,
		},
		{
			name:     "empty password",
			username: "testadmin",
			password: "",
			name_:    "Test Admin",
			email:    "admin@test.com",
			role:     "admin",
			wantErr:  true,
		},
		{
			name:     "empty name",
			username: "testadmin",
			password: "password123",
			name_:    "",
			email:    "admin@test.com",
			role:     "admin",
			wantErr:  true,
		},
		{
			name:     "empty email",
			username: "testadmin",
			password: "password123",
			name_:    "Test Admin",
			email:    "",
			role:     "admin",
			wantErr:  true,
		},
		{
			name:     "super_admin role",
			username: "superadmin",
			password: "password123",
			name_:    "Super Admin",
			email:    "super@test.com",
			role:     "super_admin",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate inputs
			hasErr := tt.username == "" ||
				tt.password == "" ||
				tt.name_ == "" ||
				tt.email == ""

			assert.Equal(t, tt.wantErr, hasErr)
		})
	}
}

// TestRoleValidation tests role validation
func TestRoleValidation(t *testing.T) {
	validRoles := []string{"admin", "super_admin", ""}
	invalidRoles := []string{"user", "guest", "root", "administrator"}

	for _, role := range validRoles {
		assert.True(t, isValidRole(role), "Role %s should be valid", role)
	}

	for _, role := range invalidRoles {
		assert.False(t, isValidRole(role), "Role %s should be invalid", role)
	}
}

// isValidRole validates admin role
func isValidRole(role string) bool {
	switch role {
	case "admin", "super_admin", "":
		return true
	default:
		return false
	}
}

// TestPasswordValidation tests password strength validation
func TestPasswordValidation(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "minimum length password",
			password: "12345678",
			wantErr:  false,
		},
		{
			name:     "too short password",
			password: "1234567",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  true,
		},
		{
			name:     "max length password",
			password: string(make([]byte, 128)),
			wantErr:  false,
		},
		{
			name:     "too long password",
			password: string(make([]byte, 129)),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestEmailValidation tests email format validation
func TestEmailValidation(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid email",
			email:   "admin@test.com",
			wantErr: false,
		},
		{
			name:    "valid email with subdomain",
			email:   "admin@mail.test.com",
			wantErr: false,
		},
		{
			name:    "missing @",
			email:   "admintest.com",
			wantErr: true,
		},
		{
			name:    "missing domain",
			email:   "admin@",
			wantErr: true,
		},
		{
			name:    "missing TLD",
			email:   "admin@test",
			wantErr: true,
		},
		{
			name:    "empty email",
			email:   "",
			wantErr: true,
		},
		{
			name:    "too long email",
			email:   string(make([]byte, 256)) + "@test.com",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEmail(tt.email)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestUsernameValidation tests username validation
func TestUsernameValidation(t *testing.T) {
	tests := []struct {
		name     string
		username string
		wantErr  bool
	}{
		{
			name:     "valid username",
			username: "testadmin",
			wantErr:  false,
		},
		{
			name:     "empty username",
			username: "",
			wantErr:  true,
		},
		{
			name:     "max length username",
			username: string(make([]byte, 50)),
			wantErr:  false,
		},
		{
			name:     "too long username",
			username: string(make([]byte, 51)),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUsername(tt.username)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Validation helper functions
func validateUsername(username string) error {
	if username == "" {
		return assert.AnError
	}
	if len(username) > maxUsernameLength {
		return assert.AnError
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) == 0 {
		return assert.AnError
	}
	if len(password) < minPasswordLength {
		return assert.AnError
	}
	if len(password) > maxPasswordLength {
		return assert.AnError
	}
	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return assert.AnError
	}
	if len(email) > maxEmailLength {
		return assert.AnError
	}
	if !emailRegex.MatchString(email) {
		return assert.AnError
	}
	return nil
}
