// Package tenant defines the tenant domain entity
package tenant

import (
	"errors"
	"time"
)

// Errors for tenant validation
var (
	ErrInvalidStatus = errors.New("invalid tenant status")
)

// Status represents the tenant status
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusSuspended Status = "suspended"
)

// IsValid checks if the status is valid
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusSuspended:
		return true
	default:
		return false
	}
}

// Tenant represents a tenant in the multi-tenant system
type Tenant struct {
	ID        string    `json:"id" db:"ID"`
	Name      string    `json:"name" db:"Name"`
	Status    Status    `json:"status" db:"Status"`
	CreatedAt time.Time `json:"createdAt" db:"CreatedAt"`
	UpdatedAt time.Time `json:"updatedAt" db:"UpdatedAt"`
}

// IsActive returns true if the tenant is active
func (t *Tenant) IsActive() bool {
	return t.Status == StatusActive
}

// Validate validates the tenant fields
func (t *Tenant) Validate() error {
	if t.Name == "" {
		return errors.New("tenant name is required")
	}
	if t.Status != "" && !t.Status.IsValid() {
		return ErrInvalidStatus
	}
	return nil
}
