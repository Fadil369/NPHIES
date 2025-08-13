package models

import "time"

// Authentication models

type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"user@nphies.sa"`
	Password string `json:"password" binding:"required" example:"password123"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	TokenType    string `json:"token_type" example:"Bearer"`
	ExpiresIn    int    `json:"expires_in" example:"3600"`
	RefreshToken string `json:"refresh_token,omitempty" example:"refresh_token_here"`
	Scope        string `json:"scope" example:"read write"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Error response model

type ErrorResponse struct {
	Error   string `json:"error" example:"invalid_request"`
	Message string `json:"message" example:"The request is missing a required parameter"`
}

// Eligibility models

type EligibilityRequest struct {
	MemberID     string `json:"member_id" binding:"required" example:"12345678901"`
	ProviderID   string `json:"provider_id" binding:"required" example:"PRV001"`
	ServiceDate  string `json:"service_date" binding:"required" example:"2025-08-13"`
	ServiceCodes []string `json:"service_codes,omitempty" example:"99213,99214"`
}

type EligibilityResponse struct {
	MemberID       string                 `json:"member_id" example:"12345678901"`
	Eligible       bool                   `json:"eligible" example:"true"`
	CoverageStatus string                 `json:"coverage_status" example:"active"`
	EffectiveDate  string                 `json:"effective_date" example:"2025-01-01"`
	ExpirationDate string                 `json:"expiration_date" example:"2025-12-31"`
	Benefits       []BenefitInformation   `json:"benefits"`
	Limitations    []CoverageLimitation   `json:"limitations"`
	Messages       []ResponseMessage      `json:"messages"`
}

type CoverageResponse struct {
	MemberID   string                 `json:"member_id" example:"12345678901"`
	Coverages  []CoverageInformation  `json:"coverages"`
	Messages   []ResponseMessage      `json:"messages"`
}

type BenefitInformation struct {
	ServiceCategory string  `json:"service_category" example:"medical"`
	CopayAmount     float64 `json:"copay_amount" example:"25.00"`
	CoinsuranceRate float64 `json:"coinsurance_rate" example:"0.20"`
	DeductibleMet   bool    `json:"deductible_met" example:"false"`
	RemainingDeductible float64 `json:"remaining_deductible" example:"500.00"`
}

type CoverageLimitation struct {
	ServiceCategory string `json:"service_category" example:"dental"`
	LimitationType  string `json:"limitation_type" example:"annual_maximum"`
	LimitValue      string `json:"limit_value" example:"2000.00"`
	UsedAmount      string `json:"used_amount" example:"450.00"`
	RemainingAmount string `json:"remaining_amount" example:"1550.00"`
}

type CoverageInformation struct {
	PayerID        string    `json:"payer_id" example:"PAY001"`
	PayerName      string    `json:"payer_name" example:"Saudi Insurance Company"`
	PolicyNumber   string    `json:"policy_number" example:"POL123456"`
	GroupNumber    string    `json:"group_number" example:"GRP789"`
	EffectiveDate  string    `json:"effective_date" example:"2025-01-01"`
	ExpirationDate string    `json:"expiration_date" example:"2025-12-31"`
	Status         string    `json:"status" example:"active"`
}

// Claims models

type ClaimSubmission struct {
	ProviderID      string           `json:"provider_id" binding:"required" example:"PRV001"`
	MemberID        string           `json:"member_id" binding:"required" example:"12345678901"`
	ServiceDate     string           `json:"service_date" binding:"required" example:"2025-08-13"`
	ClaimLines      []ClaimLine      `json:"claim_lines" binding:"required"`
	DiagnosisCodes  []DiagnosisCode  `json:"diagnosis_codes" binding:"required"`
	Attachments     []Attachment     `json:"attachments,omitempty"`
	IdempotencyKey  string           `json:"idempotency_key,omitempty" example:"uuid-v4-here"`
}

type ClaimLine struct {
	ServiceCode     string  `json:"service_code" binding:"required" example:"99213"`
	ServiceDate     string  `json:"service_date" binding:"required" example:"2025-08-13"`
	Units           int     `json:"units" example:"1"`
	ChargedAmount   float64 `json:"charged_amount" binding:"required" example:"150.00"`
	Modifiers       []string `json:"modifiers,omitempty" example:"25,59"`
	PlaceOfService  string  `json:"place_of_service" example:"11"`
}

type DiagnosisCode struct {
	Code     string `json:"code" binding:"required" example:"Z00.00"`
	CodeType string `json:"code_type" example:"ICD-10"`
	Primary  bool   `json:"primary" example:"true"`
}

type Attachment struct {
	Type        string `json:"type" example:"medical_record"`
	Description string `json:"description" example:"Patient X-ray"`
	ContentType string `json:"content_type" example:"image/jpeg"`
	Data        string `json:"data" example:"base64_encoded_data"`
}

type ClaimSubmissionResponse struct {
	ClaimID        string    `json:"claim_id" example:"CLM123456"`
	Status         string    `json:"status" example:"submitted"`
	SubmissionDate time.Time `json:"submission_date" example:"2025-08-13T10:30:00Z"`
	TrackingNumber string    `json:"tracking_number" example:"TRK789012"`
	Messages       []ResponseMessage `json:"messages"`
}

type ClaimStatusResponse struct {
	ClaimID           string               `json:"claim_id" example:"CLM123456"`
	Status            string               `json:"status" example:"approved"`
	StatusDate        time.Time            `json:"status_date" example:"2025-08-14T15:45:00Z"`
	ProcessingStages  []ProcessingStage    `json:"processing_stages"`
	AdjudicationResults []AdjudicationResult `json:"adjudication_results,omitempty"`
	Messages          []ResponseMessage    `json:"messages"`
}

type ClaimReprocessResponse struct {
	ClaimID         string    `json:"claim_id" example:"CLM123456"`
	ReprocessStatus string    `json:"reprocess_status" example:"queued"`
	RequestDate     time.Time `json:"request_date" example:"2025-08-13T16:00:00Z"`
	EstimatedDate   time.Time `json:"estimated_completion_date" example:"2025-08-14T10:00:00Z"`
}

type ProcessingStage struct {
	Stage       string    `json:"stage" example:"initial_validation"`
	Status      string    `json:"status" example:"completed"`
	StartDate   time.Time `json:"start_date" example:"2025-08-13T10:30:00Z"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"2025-08-13T10:31:00Z"`
	Messages    []ResponseMessage `json:"messages"`
}

type AdjudicationResult struct {
	LineNumber      int     `json:"line_number" example:"1"`
	ServiceCode     string  `json:"service_code" example:"99213"`
	ChargedAmount   float64 `json:"charged_amount" example:"150.00"`
	AllowedAmount   float64 `json:"allowed_amount" example:"120.00"`
	DeductibleAmount float64 `json:"deductible_amount" example:"25.00"`
	CopayAmount     float64 `json:"copay_amount" example:"25.00"`
	CoinsuranceAmount float64 `json:"coinsurance_amount" example:"19.00"`
	PaidAmount      float64 `json:"paid_amount" example:"76.00"`
	Status          string  `json:"status" example:"approved"`
	ReasonCodes     []string `json:"reason_codes,omitempty" example:"01,45"`
}

// Terminology models

type CodeSystem struct {
	ID          string `json:"id" example:"icd-10"`
	Name        string `json:"name" example:"ICD-10-CM"`
	Version     string `json:"version" example:"2025"`
	Description string `json:"description" example:"International Classification of Diseases, 10th Revision, Clinical Modification"`
	Publisher   string `json:"publisher" example:"WHO"`
	Status      string `json:"status" example:"active"`
}

type CodeLookupResponse struct {
	Code        string `json:"code" example:"Z00.00"`
	Display     string `json:"display" example:"Encounter for general adult medical examination without abnormal findings"`
	System      string `json:"system" example:"icd-10"`
	Version     string `json:"version" example:"2025"`
	Active      bool   `json:"active" example:"true"`
	Description string `json:"description,omitempty"`
}

type CodeValidationRequest struct {
	Code    string `json:"code" binding:"required" example:"Z00.00"`
	Display string `json:"display,omitempty" example:"General medical examination"`
}

type CodeValidationResponse struct {
	Valid   bool   `json:"valid" example:"true"`
	Code    string `json:"code" example:"Z00.00"`
	Display string `json:"display" example:"Encounter for general adult medical examination without abnormal findings"`
	Message string `json:"message,omitempty"`
}

// Administrative models

type SystemStats struct {
	Service      string                 `json:"service" example:"api-gateway"`
	Version      string                 `json:"version" example:"1.0.0"`
	Uptime       string                 `json:"uptime" example:"72h30m45s"`
	RequestStats RequestStatistics      `json:"request_stats"`
	Resources    ResourceUtilization    `json:"resources"`
	Dependencies DependencyStatus       `json:"dependencies"`
}

type RequestStatistics struct {
	TotalRequests       int64   `json:"total_requests" example:"125000"`
	RequestsPerSecond   float64 `json:"requests_per_second" example:"45.2"`
	AverageResponseTime float64 `json:"average_response_time_ms" example:"120.5"`
	ErrorRate           float64 `json:"error_rate" example:"0.02"`
}

type ResourceUtilization struct {
	CPUUsage    float64 `json:"cpu_usage_percent" example:"35.2"`
	MemoryUsage float64 `json:"memory_usage_percent" example:"68.5"`
	DiskUsage   float64 `json:"disk_usage_percent" example:"25.1"`
}

type DependencyStatus struct {
	Database bool `json:"database" example:"true"`
	Redis    bool `json:"redis" example:"true"`
	Kafka    bool `json:"kafka" example:"true"`
}

type AuditLogResponse struct {
	Logs       []AuditLogEntry `json:"logs"`
	TotalCount int             `json:"total_count" example:"1500"`
	Page       int             `json:"page" example:"1"`
	PageSize   int             `json:"page_size" example:"100"`
}

type AuditLogEntry struct {
	EventID     string                 `json:"event_id" example:"evt-123456"`
	EventType   string                 `json:"event_type" example:"fhir.patient.read"`
	UserID      string                 `json:"user_id" example:"user@nphies.sa"`
	ClientIP    string                 `json:"client_ip" example:"192.168.1.100"`
	Timestamp   time.Time              `json:"timestamp" example:"2025-08-13T10:30:00Z"`
	Service     string                 `json:"service" example:"api-gateway"`
	Resource    string                 `json:"resource,omitempty" example:"Patient/123"`
	Action      string                 `json:"action" example:"read"`
	Status      string                 `json:"status" example:"success"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// Common models

type ResponseMessage struct {
	Type    string `json:"type" example:"information"`
	Code    string `json:"code" example:"INFO001"`
	Message string `json:"message" example:"Request processed successfully"`
	Details string `json:"details,omitempty"`
}