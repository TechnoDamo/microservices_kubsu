-- init-db.sql
-- Create the main database
CREATE DATABASE kubsu_project_db;

-- Connect to it (the init script runs as 'postgres' superuser)
\c kubsu_project_db;

-- Create schemas for your microservices
CREATE SCHEMA IF NOT EXISTS obauth;
CREATE SCHEMA IF NOT EXISTS obprofiles;
CREATE SCHEMA IF NOT EXISTS obreports;
CREATE SCHEMA IF NOT EXISTS obnotifications;

-- Create the registered_client table in obauth schema
CREATE TABLE IF NOT EXISTS obauth.registered_client (
    id SERIAL PRIMARY KEY,  -- Auto-incrementing integer ID
    login VARCHAR(255) UNIQUE NOT NULL,  -- Unique username
    password VARCHAR(255) NOT NULL,  -- Store password hash here (not plain text!)
    is_active BOOLEAN DEFAULT TRUE,  -- User account status (active/inactive)
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP  -- Auto-set on insert
);