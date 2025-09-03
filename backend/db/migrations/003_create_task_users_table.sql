-- +goose Up
-- Migration: 003_create_task_users_table.sql
-- Description: Create junction table for many-to-many relationship between tasks and users

CREATE TABLE IF NOT EXISTS task_users (
    task_id TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (task_id, user_id),
    FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for faster joins
CREATE INDEX IF NOT EXISTS idx_task_users_task_id ON task_users(task_id);
CREATE INDEX IF NOT EXISTS idx_task_users_user_id ON task_users(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_task_users_user_id;
DROP INDEX IF EXISTS idx_task_users_task_id;
DROP TABLE IF EXISTS task_users;
