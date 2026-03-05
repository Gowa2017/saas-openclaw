package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/gowa/saas-openclaw/backend/internal/domain/user"
)

// TenantUserRepository handles database operations for tenant users
type TenantUserRepository struct {
	db *sqlx.DB
}

// NewTenantUserRepository creates a new TenantUserRepository
func NewTenantUserRepository(db *sqlx.DB) *TenantUserRepository {
	return &TenantUserRepository{db: db}
}

// Create inserts a new tenant user into the database
func (r *TenantUserRepository) Create(u *user.TenantUser) error {
	// Validate input
	if err := u.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	if u.ID == "" {
		u.ID = uuid.New().String()
	}

	// Set default role if not specified
	if u.Role == "" {
		u.Role = user.RoleTenantUser
	}

	query := `
		INSERT INTO tenant_users ("ID", "TenantID", "Name", "Email", "Role", "CreatedAt", "UpdatedAt")
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING "CreatedAt", "UpdatedAt"
	`

	err := r.db.QueryRowx(query, u.ID, u.TenantID, u.Name, u.Email, u.Role).Scan(&u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create tenant user: %w", err)
	}

	return nil
}

// GetByID retrieves a tenant user by ID
func (r *TenantUserRepository) GetByID(id string) (*user.TenantUser, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrValidation)
	}

	var u user.TenantUser
	query := `SELECT "ID", "TenantID", "Name", "Email", "Role", "CreatedAt", "UpdatedAt" FROM tenant_users WHERE "ID" = $1`

	err := r.db.Get(&u, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get tenant user: %w", err)
	}

	return &u, nil
}

// GetByEmail retrieves a tenant user by email
func (r *TenantUserRepository) GetByEmail(email string) (*user.TenantUser, error) {
	if email == "" {
		return nil, fmt.Errorf("%w: email is required", ErrValidation)
	}

	var u user.TenantUser
	query := `SELECT "ID", "TenantID", "Name", "Email", "Role", "CreatedAt", "UpdatedAt" FROM tenant_users WHERE "Email" = $1`

	err := r.db.Get(&u, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get tenant user by email: %w", err)
	}

	return &u, nil
}

// GetByTenantID retrieves all users for a tenant with pagination
func (r *TenantUserRepository) GetByTenantID(tenantID string, limit, offset int) ([]*user.TenantUser, error) {
	if tenantID == "" {
		return nil, fmt.Errorf("%w: tenant id is required", ErrValidation)
	}

	// Apply default limit if not specified
	if limit <= 0 {
		limit = 100
	}

	var users []*user.TenantUser
	query := `SELECT "ID", "TenantID", "Name", "Email", "Role", "CreatedAt", "UpdatedAt" FROM tenant_users WHERE "TenantID" = $1 ORDER BY "CreatedAt" DESC LIMIT $2 OFFSET $3`

	err := r.db.Select(&users, query, tenantID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant users: %w", err)
	}

	return users, nil
}

// Update updates an existing tenant user
func (r *TenantUserRepository) Update(u *user.TenantUser) error {
	// Validate input
	if err := u.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	if u.ID == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}

	query := `
		UPDATE tenant_users
		SET "Name" = $1, "Email" = $2, "Role" = $3, "UpdatedAt" = NOW()
		WHERE "ID" = $4
		RETURNING "UpdatedAt"
	`

	err := r.db.QueryRowx(query, u.Name, u.Email, u.Role, u.ID).Scan(&u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update tenant user: %w", err)
	}

	return nil
}

// Delete removes a tenant user by ID
func (r *TenantUserRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}

	query := `DELETE FROM tenant_users WHERE "ID" = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant user: %w", err)
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
