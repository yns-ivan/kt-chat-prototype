-- KT Chat Database Initialization Script
-- This script sets up the database and initial data for the KT Chat prototype

-- Create database if it doesn't exist
CREATE DATABASE IF NOT EXISTS ktchat CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the database
USE ktchat;

-- Create user if it doesn't exist (for local development)
-- In production, this should be handled by the deployment process
CREATE USER IF NOT EXISTS 'ktchat'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON ktchat.* TO 'ktchat'@'%';
FLUSH PRIVILEGES;

-- Note: Tables will be created automatically by GORM auto-migration
-- This script is mainly for database and user setup

-- Optional: Insert some initial data for testing
-- INSERT INTO users (id, username, email, cognito_id, created_at, updated_at) VALUES
-- ('user-123', 'admin', 'admin@example.com', 'cognito-admin-id', NOW(), NOW());

-- INSERT INTO chat_rooms (id, name, description, created_by, is_private, created_at, updated_at) VALUES
-- ('room-123', 'General', 'General discussion room', 'user-123', false, NOW(), NOW());

-- INSERT INTO room_participants (id, room_id, user_id, joined_at) VALUES
-- ('participant-123', 'room-123', 'user-123', NOW()); 