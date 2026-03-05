package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"

	"github.com/gowa/saas-openclaw/backend/internal/domain/user"
)

// AdminUserRepository handles database operations for admin users
type AdminUserRepository struct {
	db *sqlx.DB
}

// NewAdminUserRepository creates a new AdminUserRepository
func NewAdminUserRepository(db *sqlx.DB) *AdminUserRepository {
	return &AdminUserRepository{db: db}
}

// Create inserts a new admin user into the database
// The password will be hashed using bcrypt before storage
func (r *AdminUserRepository) Create(u *user.AdminUser, password string) error {
	// Validate input
	if err := u.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	// Validate password
	if password == "" {
		return fmt.Errorf("%w: password is required", ErrValidation)
	}

	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	// Set default role if not specified
	if u.Role == "" {
		u.Role = user.AdminRoleAdmin
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.PasswordHash = string(hashedPassword)

	query := `
		INSERT INTO admin_users ("ID", "Username", "PasswordHash", "Name", "Email", "Role", "CreatedAt", "UpdatedAt")
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING "CreatedAt", "UpdatedAt"
	`

	err = r.db.QueryRowx(query, u.ID, u.Username, u.PasswordHash, u.Name, u.Email, u.Role).Scan(&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	return nil
}

// GetByID retrieves an admin user by ID
func (r *AdminUserRepository) GetByID(id string) (*user.AdminUser, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrValidation)
	}

	var u user.AdminUser
	query := `SELECT "ID", "Username", "PasswordHash", "Name", "Email", "Role", "CreatedAt", "UpdatedAt" FROM admin_users WHERE "ID" = $1`

	err := r.db.Get(&u, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get admin user: %w", err)
	}

	return &u, nil
}

// GetByUsername retrieves an admin user by username
func (r *AdminUserRepository) GetByUsername(username string) (*user.AdminUser, error) {
	if username == "" {
		return nil, fmt.Errorf("%w: username is required", ErrValidation)
	}

	var u user.AdminUser
	query := `SELECT "ID", "Username", "PasswordHash", "Name", "Email", "Role", "CreatedAt", "UpdatedAt" FROM admin_users WHERE "Username" = $1`

	err := r.db.Get(&u, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get admin user by username: %w", err)
	}

	return &u, nil
}

// GetByEmail retrieves an admin user by email
func (r *AdminUserRepository) GetByEmail(email string) (*user.AdminUser, error) {
	if email == "" {
		return nil, fmt.Errorf("%w: email is required", ErrValidation)
	}

	var u user.AdminUser
	query := `SELECT "ID", "Username", "PasswordHash", "Name", "Email", "Role", "CreatedAt", "UpdatedAt" FROM admin_users WHERE "Email" = $1`

	err := r.db.Get(&u, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get admin user by email: %w", err)
	}

	return &u, nil
}

// List retrieves all admin users with pagination
func (r *AdminUserRepository) List(limit, offset int) ([]*user.AdminUser, error) {
	// Apply default limit if not specified
	if limit <= 0 {
		limit = 100
	}

	var users []*user.AdminUser
	query := `SELECT "ID", "Username", "PasswordHash", "Name", "Email", "Role", "CreatedAt", "UpdatedAt" FROM admin_users ORDER BY "CreatedAt" DESC LIMIT $1 OFFSET $2`

	err := r.db.Select(&users, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list admin users: %w", err)
	}

	return users, nil
}

// Update updates an existing admin user
func (r *AdminUserRepository) Update(u *user.AdminUser) error {
	// Validate input
	if err := u.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	if u.ID == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}

	query := `
		UPDATE admin_users
		SET "Name" = $1, "Email" = $2, "Role" = $3, "UpdatedAt" = NOW()
		WHERE "ID" = $4
		RETURNING "UpdatedAt"
	`

	err := r.db.QueryRowx(query, u.Name, u.Email, u.Role, u.ID).Scan(&u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update admin user: %w", err)
	}

	return nil
}

// UpdatePassword updates an admin user's password
func (r *AdminUserRepository) UpdatePassword(id string, newPassword string) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}

	// Validate new password
	if newPassword == "" {
		return fmt.Errorf("%w: new password is required", ErrValidation)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	query := `UPDATE admin_users SET "PasswordHash" = $1, "UpdatedAt" = NOW() WHERE "ID" = $2`

	result, err := r.db.Exec(query, string(hashedPassword), id)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

// VerifyPassword verifies an admin user's password
func (r *AdminUserRepository) VerifyPassword(u *user.AdminUser, password string) bool {
	if u == nil || u.PasswordHash == "" || password == "" {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// Delete removes an admin user by ID
func (r *AdminUserRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}

	query := `DELETE FROM admin_users WHERE "ID" = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete admin user: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
