-- Eligibility Service Database Schema
-- Connect to eligibility database
\c eligibility;

-- Create UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create members table
CREATE TABLE IF NOT EXISTS members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identifier VARCHAR(255) UNIQUE NOT NULL, -- National ID
    name JSONB NOT NULL, -- FHIR HumanName structure
    birth_date DATE NOT NULL,
    gender VARCHAR(20) NOT NULL,
    contact_info JSONB,
    address JSONB,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create coverage table
CREATE TABLE IF NOT EXISTS coverage (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    member_id VARCHAR(255) NOT NULL,
    payer_id VARCHAR(255) NOT NULL,
    policy_number VARCHAR(255),
    group_number VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active',
    type VARCHAR(50) NOT NULL, -- medical, dental, vision, etc.
    effective_date DATE NOT NULL,
    expiration_date DATE,
    benefit_details JSONB,
    cost_sharing JSONB,
    network VARCHAR(255),
    prior_auth_rules JSONB,
    limitations JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (member_id) REFERENCES members(identifier)
);

-- Create providers table
CREATE TABLE IF NOT EXISTS providers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identifier VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- hospital, clinic, pharmacy, etc.
    organization_id UUID,
    specialties TEXT[],
    contact_info JSONB,
    address JSONB,
    network_affiliations TEXT[],
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create eligibility_checks table for audit and caching
CREATE TABLE IF NOT EXISTS eligibility_checks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id VARCHAR(255) NOT NULL,
    member_id VARCHAR(255) NOT NULL,
    provider_id VARCHAR(255),
    service_date DATE NOT NULL,
    service_codes TEXT[],
    response_data JSONB NOT NULL,
    cache_key VARCHAR(500),
    processing_time_ms INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create benefit_utilization table
CREATE TABLE IF NOT EXISTS benefit_utilization (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    member_id VARCHAR(255) NOT NULL,
    coverage_id UUID NOT NULL,
    service_category VARCHAR(100) NOT NULL,
    period VARCHAR(20) NOT NULL, -- annual, monthly, lifetime
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    limit_amount DECIMAL(10,2),
    used_amount DECIMAL(10,2) DEFAULT 0,
    remaining_amount DECIMAL(10,2),
    transaction_count INTEGER DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (coverage_id) REFERENCES coverage(id),
    UNIQUE(member_id, coverage_id, service_category, period_start)
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_members_identifier ON members(identifier);
CREATE INDEX IF NOT EXISTS idx_members_status ON members(status);

CREATE INDEX IF NOT EXISTS idx_coverage_member_id ON coverage(member_id);
CREATE INDEX IF NOT EXISTS idx_coverage_payer_id ON coverage(payer_id);
CREATE INDEX IF NOT EXISTS idx_coverage_status ON coverage(status);
CREATE INDEX IF NOT EXISTS idx_coverage_effective_date ON coverage(effective_date);
CREATE INDEX IF NOT EXISTS idx_coverage_member_effective ON coverage(member_id, effective_date);

CREATE INDEX IF NOT EXISTS idx_providers_identifier ON providers(identifier);
CREATE INDEX IF NOT EXISTS idx_providers_type ON providers(type);
CREATE INDEX IF NOT EXISTS idx_providers_status ON providers(status);

CREATE INDEX IF NOT EXISTS idx_eligibility_checks_member_id ON eligibility_checks(member_id);
CREATE INDEX IF NOT EXISTS idx_eligibility_checks_service_date ON eligibility_checks(service_date);
CREATE INDEX IF NOT EXISTS idx_eligibility_checks_cache_key ON eligibility_checks(cache_key);
CREATE INDEX IF NOT EXISTS idx_eligibility_checks_created_at ON eligibility_checks(created_at);

CREATE INDEX IF NOT EXISTS idx_benefit_utilization_member_id ON benefit_utilization(member_id);
CREATE INDEX IF NOT EXISTS idx_benefit_utilization_coverage_id ON benefit_utilization(coverage_id);
CREATE INDEX IF NOT EXISTS idx_benefit_utilization_period ON benefit_utilization(period_start, period_end);

-- Grant permissions to nphies user
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO nphies;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO nphies;

-- Create triggers for updated_at
CREATE TRIGGER update_members_updated_at BEFORE UPDATE ON members
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_coverage_updated_at BEFORE UPDATE ON coverage
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_providers_updated_at BEFORE UPDATE ON providers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert sample data for testing
INSERT INTO members (identifier, name, birth_date, gender, contact_info, address, status) VALUES
('1234567890', '{"use": "official", "family": "العتيبي", "given": ["أحمد", "محمد"]}', '1985-03-15', 'male', 
 '{"phone": "+966501234567", "email": "ahmed.alotaibi@example.com"}', 
 '{"line": ["123 King Fahd Road"], "city": "Riyadh", "state": "Riyadh Province", "postalCode": "12345", "country": "SA"}', 'active'),
('0987654321', '{"use": "official", "family": "الزهراني", "given": ["فاطمة", "علي"]}', '1990-07-22', 'female',
 '{"phone": "+966507654321", "email": "fatima.alzahrani@example.com"}',
 '{"line": ["456 Prince Mohammed Road"], "city": "Jeddah", "state": "Makkah Province", "postalCode": "54321", "country": "SA"}', 'active')
ON CONFLICT (identifier) DO NOTHING;

INSERT INTO providers (identifier, name, type, specialties, contact_info, address, network_affiliations, status) VALUES
('PRV001', 'King Faisal Specialist Hospital', 'hospital', ARRAY['cardiology', 'oncology', 'surgery'], 
 '{"phone": "+966114647272", "email": "info@kfshrc.edu.sa"}',
 '{"line": ["Maather Street"], "city": "Riyadh", "state": "Riyadh Province", "postalCode": "11211", "country": "SA"}',
 ARRAY['network_a', 'network_b'], 'active'),
('PRV002', 'Al-Jedaani Hospital', 'hospital', ARRAY['emergency', 'internal_medicine'],
 '{"phone": "+966126651234", "email": "info@aljedaani.com"}',
 '{"line": ["Prince Sultan Street"], "city": "Jeddah", "state": "Makkah Province", "postalCode": "21589", "country": "SA"}',
 ARRAY['network_a'], 'active')
ON CONFLICT (identifier) DO NOTHING;

INSERT INTO coverage (member_id, payer_id, policy_number, group_number, status, type, effective_date, expiration_date, 
                     benefit_details, cost_sharing, network, prior_auth_rules, limitations) VALUES
('1234567890', 'PAY001', 'POL123456', 'GRP789', 'active', 'medical', '2025-01-01', '2025-12-31',
 '{"medical": {"annual_limit": 100000, "copay": 25, "coinsurance": 0.2}}',
 '{"deductible": 500, "out_of_pocket_max": 2000}', 'network_a',
 '{"high_cost_services": ["surgery", "imaging"]}',
 '{"annual_maximum": 100000, "visit_limits": {"specialist": 12, "therapy": 20}}'),
('0987654321', 'PAY002', 'POL654321', 'GRP456', 'active', 'medical', '2025-01-01', '2025-12-31',
 '{"medical": {"annual_limit": 75000, "copay": 50, "coinsurance": 0.25}}',
 '{"deductible": 1000, "out_of_pocket_max": 3000}', 'network_b',
 '{}', '{"annual_maximum": 75000}')
ON CONFLICT DO NOTHING;

INSERT INTO benefit_utilization (member_id, coverage_id, service_category, period, period_start, period_end, 
                               limit_amount, used_amount, remaining_amount, transaction_count) VALUES
('1234567890', (SELECT id FROM coverage WHERE member_id = '1234567890' LIMIT 1), 'medical', 'annual', 
 '2025-01-01', '2025-12-31', 100000.00, 15000.00, 85000.00, 12),
('0987654321', (SELECT id FROM coverage WHERE member_id = '0987654321' LIMIT 1), 'medical', 'annual',
 '2025-01-01', '2025-12-31', 75000.00, 8500.00, 66500.00, 8)
ON CONFLICT DO NOTHING;