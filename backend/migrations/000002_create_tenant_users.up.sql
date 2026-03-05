-- +migrate Up
-- Create tenant_users table for multi-tenant user management
CREATE TABLE tenant_users (
    "ID" VARCHAR(36) PRIMARY KEY,
    "TenantID" VARCHAR(36) NOT NULL REFERENCES tenants("ID") ON DELETE CASCADE,
    "Name" VARCHAR(255) NOT NULL,
    "Email" VARCHAR(255) NOT NULL,
    "Role" VARCHAR(50) NOT NULL DEFAULT 'user',
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on TenantID for tenant-based queries
CREATE INDEX idx_tenant_users_tenant_id ON tenant_users ("TenantID");

-- Create unique index on Email for login
CREATE UNIQUE INDEX idx_tenant_users_email ON tenant_users ("Email");

-- Create index on Role for role-based queries
CREATE INDEX idx_tenant_users_role ON tenant_users ("Role");

-- Comment on table
COMMENT ON TABLE tenant_users IS 'Tenant users table for multi-tenant user management';
COMMENT ON COLUMN tenant_users."ID" IS 'UUID primary key';
COMMENT ON COLUMN tenant_users."TenantID" IS 'Foreign key to tenants table';
COMMENT ON COLUMN tenant_users."Name" IS 'User display name';
COMMENT ON COLUMN tenant_users."Email" IS 'User email address (unique)';
COMMENT ON COLUMN tenant_users."Role" IS 'User role: user, admin';
COMMENT ON COLUMN tenant_users."CreatedAt" IS 'Record creation timestamp';
COMMENT ON COLUMN tenant_users."UpdatedAt" IS 'Record update timestamp';
