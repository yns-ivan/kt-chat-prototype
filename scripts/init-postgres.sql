-- KT Chat Database Initialization Script for PostgreSQL
-- This script sets up the database and initial data for the KT Chat prototype

-- Create database if it doesn't exist (PostgreSQL doesn't have CREATE DATABASE IF NOT EXISTS)
-- The database is created by the POSTGRES_DB environment variable

-- Enable UUID extension for better UUID support
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable pgcrypto extension for cryptographic functions (if needed)
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create user if it doesn't exist (for local development)
-- In production, this should be handled by the deployment process
DO $$
BEGIN
    IF NOT EXISTS (SELECT FROM pg_catalog.pg_roles WHERE rolname = 'ktchat') THEN
        CREATE USER ktchat WITH PASSWORD 'password';
    END IF;
END
$$;

-- Grant privileges
GRANT ALL PRIVILEGES ON DATABASE ktchat TO ktchat;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ktchat;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO ktchat;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA public TO ktchat;

-- Set default privileges for future objects
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ktchat;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON SEQUENCES TO ktchat;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON FUNCTIONS TO ktchat;

-- Note: Tables will be created automatically by GORM auto-migration
-- This script is mainly for database setup and user permissions

-- Optional: Insert some initial data for testing
-- INSERT INTO users (id, username, email, cognito_id, created_at, updated_at) VALUES
-- (gen_random_uuid(), 'admin', 'admin@example.com', 'cognito-admin-id', NOW(), NOW());

-- INSERT INTO chat_rooms (id, name, description, created_by, is_private, created_at, updated_at) VALUES
-- (gen_random_uuid(), 'General', 'General discussion room', 'user-123', false, NOW(), NOW());

-- INSERT INTO room_participants (id, room_id, user_id, joined_at) VALUES
-- (gen_random_uuid(), 'room-123', 'user-123', NOW()); 