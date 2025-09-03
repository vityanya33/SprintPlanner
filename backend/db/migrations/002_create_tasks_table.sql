-- +goose Up
-- Migration: 002_create_tasks_table.sql
-- Description: Create tasks table based on Task model

CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    hours INTEGER NOT NULL DEFAULT 0,
    start_date DATE,
    deadline DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on dates for faster queries
CREATE INDEX IF NOT EXISTS idx_tasks_start_date ON tasks(start_date);
CREATE INDEX IF NOT EXISTS idx_tasks_deadline ON tasks(deadline);
CREATE INDEX IF NOT EXISTS idx_tasks_created_at ON tasks(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_tasks_created_at;
DROP INDEX IF EXISTS idx_tasks_deadline;
DROP INDEX IF EXISTS idx_tasks_start_date;
DROP TABLE IF EXISTS tasks;
