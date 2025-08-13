package models

import "time"

// Coverage represents a FHIR Coverage resource with additional business logic
type Coverage struct {
	ID               string                 `json:"id" db:"id"`
	MemberID         string                 `json:"member_id" db:"member_id"`
	PayerID          string                 `json:"payer_id" db:"payer_id"`
	PolicyNumber     string                 `json:"policy_number" db:"policy_number"`
	GroupNumber      string                 `json:"group_number" db:"group_number"`
	Status           string                 `json:"status" db:"status"` // active, cancelled, draft, entered-in-error
	Type             string                 `json:"type" db:"type"`     // medical, dental, vision, etc.
	EffectiveDate    time.Time              `json:"effective_date" db:"effective_date"`
	ExpirationDate   *time.Time             `json:"expiration_date" db:"expiration_date"`
	BenefitDetails   map[string]interface{} `json:"benefit_details" db:"benefit_details"`
	CostSharing      map[string]interface{} `json:"cost_sharing" db:"cost_sharing"`
	Network          string                 `json:"network" db:"network"`
	PriorAuthRules   map[string]interface{} `json:"prior_auth_rules" db:"prior_auth_rules"`
	Limitations      map[string]interface{} `json:"limitations" db:"limitations"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at" db:"updated_at"`
}

// Member represents a patient/member
type Member struct {
	ID           string                 `json:"id" db:"id"`
	Identifier   string                 `json:"identifier" db:"identifier"` // National ID
	Name         map[string]interface{} `json:"name" db:"name"`             // FHIR HumanName
	BirthDate    time.Time              `json:"birth_date" db:"birth_date"`
	Gender       string                 `json:"gender" db:"gender"`
	ContactInfo  map[string]interface{} `json:"contact_info" db:"contact_info"`
	Address      map[string]interface{} `json:"address" db:"address"`
	Status       string                 `json:"status" db:"status"` // active, inactive
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at" db:"updated_at"`
}

// Provider represents a healthcare provider
type Provider struct {
	ID              string                 `json:"id" db:"id"`
	Identifier      string                 `json:"identifier" db:"identifier"`
	Name            string                 `json:"name" db:"name"`
	Type            string                 `json:"type" db:"type"` // hospital, clinic, pharmacy, etc.
	OrganizationID  *string                `json:"organization_id" db:"organization_id"`
	Specialties     []string               `json:"specialties" db:"specialties"`
	ContactInfo     map[string]interface{} `json:"contact_info" db:"contact_info"`
	Address         map[string]interface{} `json:"address" db:"address"`
	NetworkAffiliations []string           `json:"network_affiliations" db:"network_affiliations"`
	Status          string                 `json:"status" db:"status"`
	CreatedAt       time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at" db:"updated_at"`
}

// EligibilityRequest represents an eligibility check request
type EligibilityRequest struct {
	RequestID    string    `json:"request_id"`
	MemberID     string    `json:"member_id" binding:"required"`
	ProviderID   string    `json:"provider_id" binding:"required"`
	ServiceDate  string    `json:"service_date" binding:"required"`
	ServiceCodes []string  `json:"service_codes,omitempty"`
	RequestedBy  string    `json:"requested_by,omitempty"`
	RequestTime  time.Time `json:"request_time"`
}

// EligibilityResponse represents the response to an eligibility check
type EligibilityResponse struct {
	RequestID      string                 `json:"request_id"`
	MemberID       string                 `json:"member_id"`
	Eligible       bool                   `json:"eligible"`
	CoverageStatus string                 `json:"coverage_status"`
	EffectiveDate  string                 `json:"effective_date"`
	ExpirationDate string                 `json:"expiration_date,omitempty"`
	Benefits       []BenefitInformation   `json:"benefits"`
	Limitations    []CoverageLimitation   `json:"limitations"`
	Messages       []ResponseMessage      `json:"messages"`
	ResponseTime   time.Time              `json:"response_time"`
	CacheHit       bool                   `json:"cache_hit"`
}

// BenefitInformation represents benefit details
type BenefitInformation struct {
	ServiceCategory     string  `json:"service_category"`
	InNetwork           bool    `json:"in_network"`
	CopayAmount         float64 `json:"copay_amount"`
	CoinsuranceRate     float64 `json:"coinsurance_rate"`
	DeductibleAmount    float64 `json:"deductible_amount"`
	DeductibleMet       bool    `json:"deductible_met"`
	RemainingDeductible float64 `json:"remaining_deductible"`
	OutOfPocketMax      float64 `json:"out_of_pocket_max"`
	RemainingOOPMax     float64 `json:"remaining_oop_max"`
	PriorAuthRequired   bool    `json:"prior_auth_required"`
	CoverageLevel       string  `json:"coverage_level"` // individual, family
}

// CoverageLimitation represents coverage limitations
type CoverageLimitation struct {
	ServiceCategory string  `json:"service_category"`
	LimitationType  string  `json:"limitation_type"` // annual_maximum, lifetime_maximum, visit_limit
	LimitValue      float64 `json:"limit_value"`
	UsedAmount      float64 `json:"used_amount"`
	RemainingAmount float64 `json:"remaining_amount"`
	Period          string  `json:"period"` // annual, lifetime, monthly
	ResetDate       string  `json:"reset_date,omitempty"`
}

// ResponseMessage represents informational messages
type ResponseMessage struct {
	Type    string `json:"type"`    // information, warning, error
	Code    string `json:"code"`    
	Message string `json:"message"` 
	Details string `json:"details,omitempty"`
}

// CoverageVerificationRequest represents a coverage verification request
type CoverageVerificationRequest struct {
	MemberID      string   `json:"member_id" binding:"required"`
	ServiceDate   string   `json:"service_date" binding:"required"`
	ServiceCodes  []string `json:"service_codes" binding:"required"`
	ProviderID    string   `json:"provider_id" binding:"required"`
	PlaceOfService string  `json:"place_of_service,omitempty"`
}

// CoverageVerificationResponse represents the verification response
type CoverageVerificationResponse struct {
	MemberID        string                    `json:"member_id"`
	VerificationID  string                    `json:"verification_id"`
	ServiceDate     string                    `json:"service_date"`
	Services        []ServiceVerification     `json:"services"`
	OverallStatus   string                    `json:"overall_status"` // covered, not_covered, partial
	AuthRequired    bool                      `json:"auth_required"`
	Messages        []ResponseMessage         `json:"messages"`
	ValidUntil      time.Time                 `json:"valid_until"`
}

// ServiceVerification represents verification for a specific service
type ServiceVerification struct {
	ServiceCode     string  `json:"service_code"`
	Status          string  `json:"status"` // covered, not_covered, requires_auth
	CoverageLevel   float64 `json:"coverage_level"` // percentage covered
	EstimatedCost   float64 `json:"estimated_cost"`
	PatientCost     float64 `json:"patient_cost"`
	AuthRequired    bool    `json:"auth_required"`
	AuthReference   string  `json:"auth_reference,omitempty"`
	ReasonCodes     []string `json:"reason_codes,omitempty"`
}

// ServiceStats represents service statistics
type ServiceStats struct {
	Service         string                 `json:"service"`
	Version         string                 `json:"version"`
	Uptime          string                 `json:"uptime"`
	RequestStats    RequestStatistics      `json:"request_stats"`
	CacheStats      CacheStatistics        `json:"cache_stats"`
	DatabaseStats   DatabaseStatistics     `json:"database_stats"`
	Dependencies    DependencyStatus       `json:"dependencies"`
}

// RequestStatistics represents request statistics
type RequestStatistics struct {
	TotalRequests       int64   `json:"total_requests"`
	RequestsPerSecond   float64 `json:"requests_per_second"`
	AverageResponseTime float64 `json:"average_response_time_ms"`
	P95ResponseTime     float64 `json:"p95_response_time_ms"`
	P99ResponseTime     float64 `json:"p99_response_time_ms"`
	ErrorRate           float64 `json:"error_rate"`
	SuccessRate         float64 `json:"success_rate"`
}

// CacheStatistics represents cache statistics
type CacheStatistics struct {
	HitRate        float64 `json:"hit_rate"`
	MissRate       float64 `json:"miss_rate"`
	TotalHits      int64   `json:"total_hits"`
	TotalMisses    int64   `json:"total_misses"`
	CachedEntries  int64   `json:"cached_entries"`
	CacheSize      string  `json:"cache_size"`
	EvictionCount  int64   `json:"eviction_count"`
}

// DatabaseStatistics represents database statistics
type DatabaseStatistics struct {
	ActiveConnections int `json:"active_connections"`
	IdleConnections   int `json:"idle_connections"`
	TotalQueries      int64 `json:"total_queries"`
	SlowQueries       int64 `json:"slow_queries"`
	AverageQueryTime  float64 `json:"average_query_time_ms"`
}

// DependencyStatus represents the status of external dependencies
type DependencyStatus struct {
	Database bool `json:"database"`
	Redis    bool `json:"redis"`
	Kafka    bool `json:"kafka"`
}

// CacheEntry represents a cached eligibility entry
type CacheEntry struct {
	Key        string                 `json:"key"`
	Data       interface{}            `json:"data"`
	CreatedAt  time.Time              `json:"created_at"`
	ExpiresAt  time.Time              `json:"expires_at"`
	AccessCount int                   `json:"access_count"`
	LastAccess time.Time              `json:"last_access"`
}