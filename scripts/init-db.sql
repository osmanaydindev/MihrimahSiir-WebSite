-- Database initialization script
-- This script runs when the PostgreSQL container starts for the first time

-- Create database if not exists (already created by POSTGRES_DB env var)
-- CREATE DATABASE poem_blog_locale;

-- Set timezone
SET timezone = 'Europe/Istanbul';

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- You can add any initial data or additional setup here
-- For example:
-- INSERT INTO admins (username, password, role_id) VALUES ('admin', 'hashed_password', 1);
