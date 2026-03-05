-- +migrate Up
-- Create admin_users table for platform administration
CREATE TABLE admin_users (
    "ID" VARCHAR(36) PRIMARY KEY,
    "Username" VARCHAR(100) NOT NULL,
    "PasswordHash" VARCHAR(255) NOT NULL,
    "Name" VARCHAR(255) NOT NULL,
    "Email" VARCHAR(255) NOT NULL,
    "Role" VARCHAR(50) NOT NULL DEFAULT 'admin',
    "CreatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "UpdatedAt" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create unique index on Username for login
CREATE UNIQUE INDEX idx_admin_users_username ON admin_users ("Username");

-- Create unique index on Email for login
CREATE UNIQUE INDEX idx_admin_users_email ON admin_users ("Email");

-- Create index on Role for role-based queries
CREATE INDEX idx_admin_users_role ON admin_users ("Role");

-- Comment on table
COMMENT ON TABLE admin_users IS 'Admin users table for platform administration';
COMMENT ON COLUMN admin_users."ID" IS 'UUID primary key';
COMMENT ON COLUMN admin_users."Username" IS 'Admin username (unique)';
COMMENT ON COLUMN admin_users."PasswordHash" IS 'Bcrypt password hash';
COMMENT ON COLUMN admin_users."Name" IS 'Admin display name';
COMMENT ON COLUMN admin_users."Email" IS 'Admin email address (unique)';
COMMENT ON COLUMN admin_users."Role" IS 'Admin role: admin, super_admin';
COMMENT ON COLUMN admin_users."CreatedAt" IS 'Record creation timestamp';
COMMENT ON COLUMN admin_users."UpdatedAt" IS 'Record update timestamp';
