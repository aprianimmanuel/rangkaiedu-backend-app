-- +goose Up
-- Migration: Add MFA and social provider support to users table
-- Description: Enhances users table for unified authentication workflow
-- Date: 2025-10-07

-- Add columns to users table for social providers
ALTER TABLE users 
ADD COLUMN google_id VARCHAR(255) UNIQUE,
ADD COLUMN facebook_id VARCHAR(255) UNIQUE;

-- Add columns for MFA support
ALTER TABLE users 
ADD COLUMN mfa_secret VARCHAR(255),
ADD COLUMN is_mfa_enabled BOOLEAN DEFAULT FALSE,
ADD COLUMN mfa_backup_codes TEXT[];

-- Update oauth_providers table to support Facebook
ALTER TABLE oauth_providers 
DROP CONSTRAINT IF EXISTS oauth_providers_provider_check;

ALTER TABLE oauth_providers 
ADD CONSTRAINT oauth_providers_provider_check 
CHECK (provider IN ('google', 'apple', 'facebook'));

-- Create MFA events table for auditing
CREATE TABLE mfa_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL, -- setup, verify, disable
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_users_google_id ON users(google_id);
CREATE INDEX IF NOT EXISTS idx_users_facebook_id ON users(facebook_id);
CREATE INDEX IF NOT EXISTS idx_users_mfa_enabled ON users(is_mfa_enabled);
CREATE INDEX IF NOT EXISTS idx_mfa_events_user_id ON mfa_events(user_id);
CREATE INDEX IF NOT EXISTS idx_mfa_events_created_at ON mfa_events(created_at);

-- +goose Down
-- Migration: Remove MFA and social provider support from users table
-- Description: Reverts users table to original schema
-- Date: 2025-10-07

-- Drop indexes
DROP INDEX IF EXISTS idx_users_google_id;
DROP INDEX IF EXISTS idx_users_facebook_id;
DROP INDEX IF EXISTS idx_users_mfa_enabled;
DROP INDEX IF EXISTS idx_mfa_events_user_id;
DROP INDEX IF EXISTS idx_mfa_events_created_at;

-- Drop MFA events table
DROP TABLE IF EXISTS mfa_events;

-- Remove constraints
ALTER TABLE oauth_providers 
DROP CONSTRAINT IF EXISTS oauth_providers_provider_check;

-- Add original constraint
ALTER TABLE oauth_providers 
ADD CONSTRAINT oauth_providers_provider_check 
CHECK (provider IN ('google', 'apple'));

-- Remove columns from users table
ALTER TABLE users 
DROP COLUMN IF EXISTS google_id,
DROP COLUMN IF EXISTS facebook_id,
DROP COLUMN IF EXISTS mfa_secret,
DROP COLUMN IF EXISTS is_mfa_enabled,
DROP COLUMN IF EXISTS mfa_backup_codes;