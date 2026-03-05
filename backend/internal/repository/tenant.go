// Package repository implements data access layer
package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/gowa/saas-openclaw/backend/internal/domain/tenant"
)

// TenantRepository handles database operations for tenants
type TenantRepository struct {
	db *sqlx.DB
}

// NewTenantRepository creates a new TenantRepository
func NewTenantRepository(db *sqlx.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

// Create inserts a new tenant into the database
func (r *TenantRepository) Create(t *tenant.Tenant) error {
	// Validate input
	if err := t.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	// Set default status if not specified
	if t.Status == "" {
		t.Status = tenant.StatusActive
	}

	query := `
		INSERT INTO tenants ("ID", "Name", "Status", "CreatedAt", "UpdatedAt")
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING "CreatedAt", "UpdatedAt"
	`

	err := r.db.QueryRowx(query, t.ID, t.Name, t.Status).Scan(&t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	return nil
}

// GetByID retrieves a tenant by ID
func (r *TenantRepository) GetByID(id string) (*tenant.Tenant, error) {
	if id == "" {
		return nil, fmt.Errorf("%w: id is required", ErrValidation)
	}

	var t tenant.Tenant
	query := `SELECT "ID", "Name", "Status", "CreatedAt", "UpdatedAt" FROM tenants WHERE "ID" = $1`

	err := r.db.Get(&t, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return &t, nil
}

// GetByName retrieves a tenant by name
func (r *TenantRepository) GetByName(name string) (*tenant.Tenant, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: name is required", ErrValidation)
	}

	var t tenant.Tenant
	query := `SELECT "ID", "Name", "Status", "CreatedAt", "UpdatedAt" FROM tenants WHERE "Name" = $1`

	err := r.db.Get(&t, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get tenant by name: %w", err)
	}

	return &t, nil
}

// List retrieves all tenants with pagination
func (r *TenantRepository) List(limit, offset int) ([]*tenant.Tenant, error) {
	var tenants []*tenant.Tenant

	// Apply default limit if not specified
	if limit <= 0 {
		limit = 100
	}

	query := `SELECT "ID", "Name", "Status", "CreatedAt", "UpdatedAt" FROM tenants ORDER BY "CreatedAt" DESC LIMIT $1 OFFSET $2`

	err := r.db.Select(&tenants, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return tenants, nil
}

// Update updates an existing tenant
func (r *TenantRepository) Update(t *tenant.Tenant) error {
	// Validate input
	if err := t.Validate(); err != nil {
		return fmt.Errorf("%w: %v", ErrValidation, err)
	}

	if t.ID == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}

	query := `
		UPDATE tenants
		SET "Name" = $1, "Status" = $2, "UpdatedAt" = NOW()
		WHERE "ID" = $3
		RETURNING "UpdatedAt"
	`

	err := r.db.QueryRowx(query, t.Name, t.Status, t.ID).Scan(&t.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound
		}
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	return nil
}

// Delete removes a tenant by ID
func (r *TenantRepository) Delete(id string) error {
	if id == "" {
		return fmt.Errorf("%w: id is required", ErrValidation)
	}

	query := `DELETE FROM tenants WHERE "ID" = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
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
