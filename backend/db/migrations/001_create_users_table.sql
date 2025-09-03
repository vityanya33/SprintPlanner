-- +goose Up
-- Migration: 001_create_users_table.sql
-- Description: Create users table based on User model

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(100) NOT NULL,
    resource INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on role for faster queries
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Create index on name for faster searches
CREATE INDEX IF NOT EXISTS idx_users_name ON users(name);

-- +goose Down
DROP INDEX IF EXISTS idx_users_name;
DROP INDEX IF EXISTS idx_users_role;
DROP TABLE IF EXISTS users;
