-- +migrate Up
-- Create tenants table for multi-tenant support
CREATE TABLE tenants (
    "ID" VARCHAR(36) PRIMARY KEY,
    "Name" VARCHAR(255) NOT NULL,
    "Status" VARCHAR(50) NOT NULL DEFAULT 'active',
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create index on status for filtering
CREATE INDEX idx_tenants_status ON tenants ("Status");

-- Create index on name for searching
CREATE INDEX idx_tenants_name ON tenants ("Name");

-- Comment on table
COMMENT ON TABLE tenants IS 'Tenant table for multi-tenant support';
COMMENT ON COLUMN tenants."ID" IS 'UUID primary key';
COMMENT ON COLUMN tenants."Name" IS 'Tenant display name';
COMMENT ON COLUMN tenants."Status" IS 'Tenant status: active, inactive, suspended';
COMMENT ON COLUMN tenants."CreatedAt" IS 'Record creation timestamp';
COMMENT ON COLUMN tenants."UpdatedAt" IS 'Record update timestamp';
