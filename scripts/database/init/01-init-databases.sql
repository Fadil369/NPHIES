-- NPHIES Platform Database Initialization Script
-- This script creates the databases and basic schema for each service

-- Create databases for each service
SELECT 'CREATE DATABASE eligibility OWNER nphies'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'eligibility')\gexec

SELECT 'CREATE DATABASE claims OWNER nphies'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'claims')\gexec

SELECT 'CREATE DATABASE terminology OWNER nphies'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'terminology')\gexec

-- Connect to main nphies database for shared tables
\c nphies;

-- Create UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create audit log table for centralized auditing
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id VARCHAR(255) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    user_id VARCHAR(255),
    client_ip INET,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    service VARCHAR(50) NOT NULL,
    resource VARCHAR(255),
    action VARCHAR(50),
    status VARCHAR(20),
    data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on audit logs for performance
CREATE INDEX IF NOT EXISTS idx_audit_logs_timestamp ON audit_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_audit_logs_event_type ON audit_logs(event_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);

-- Create system configuration table
CREATE TABLE IF NOT EXISTS system_config (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    config_key VARCHAR(255) UNIQUE NOT NULL,
    config_value JSONB NOT NULL,
    description TEXT,
    is_sensitive BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Insert default configuration
INSERT INTO system_config (config_key, config_value, description) VALUES
('rate_limit.requests_per_minute', '500', 'Default rate limit for API requests')
ON CONFLICT (config_key) DO NOTHING;

-- Create organizations table (shared across services)
CREATE TABLE IF NOT EXISTS organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identifier VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- 'payer', 'provider', 'regulator'
    status VARCHAR(20) DEFAULT 'active',
    contact_info JSONB,
    address JSONB,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create practitioners table (shared across services)
CREATE TABLE IF NOT EXISTS practitioners (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identifier VARCHAR(255) UNIQUE NOT NULL,
    name JSONB NOT NULL, -- Store FHIR HumanName structure
    qualification JSONB,
    organization_id UUID REFERENCES organizations(id),
    status VARCHAR(20) DEFAULT 'active',
    contact_info JSONB,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Grant permissions to nphies user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO nphies;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO nphies;

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_practitioners_updated_at BEFORE UPDATE ON practitioners
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();