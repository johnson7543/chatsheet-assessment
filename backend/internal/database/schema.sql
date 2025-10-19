-- Database Schema for LinkedIn Connector Application
-- This file is for reference only. GORM auto-migration will create these tables automatically.

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Linked Accounts table
CREATE TABLE IF NOT EXISTS linked_accounts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    provider TEXT NOT NULL DEFAULT 'linkedin',
    account_id TEXT NOT NULL,
    account_name TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_linked_accounts_user_id ON linked_accounts(user_id);
CREATE INDEX IF NOT EXISTS idx_linked_accounts_deleted_at ON linked_accounts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_linked_accounts_provider ON linked_accounts(provider);

-- Example queries:

-- Get all accounts for a user
-- SELECT * FROM linked_accounts WHERE user_id = ? AND deleted_at IS NULL ORDER BY created_at DESC;

-- Get user with their accounts
-- SELECT u.*, la.* FROM users u 
-- LEFT JOIN linked_accounts la ON u.id = la.user_id AND la.deleted_at IS NULL
-- WHERE u.id = ? AND u.deleted_at IS NULL;

-- Count accounts by provider
-- SELECT provider, COUNT(*) as count FROM linked_accounts 
-- WHERE deleted_at IS NULL GROUP BY provider;

