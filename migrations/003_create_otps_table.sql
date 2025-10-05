-- Migration: Create otps table for OTP storage
-- This table stores temporary OTP codes for email/phone verification
-- with expiry to prevent replay attacks

CREATE TABLE IF NOT EXISTS otps (
    id SERIAL PRIMARY KEY,
    identifier VARCHAR(255) NOT NULL,  -- email or phone number
    otp VARCHAR(6) NOT NULL,           -- 6-digit OTP code
    expiry TIMESTAMP NOT NULL,         -- OTP expiry time (e.g., 10 minutes from creation)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(identifier, otp)            -- Prevent duplicate OTPs for same identifier
);

-- Index for faster lookups by identifier and expiry
CREATE INDEX IF NOT EXISTS idx_otps_identifier_expiry ON otps (identifier, expiry);

-- Optional: Add a cleanup job or trigger to remove expired OTPs periodically