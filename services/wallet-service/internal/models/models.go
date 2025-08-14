package models

import (
	"time"
)

// WalletTransaction represents a financial transaction in the digital wallet
type WalletTransaction struct {
	ID            string    `json:"id" db:"id"`
	MemberID      string    `json:"member_id" db:"member_id"`
	TransactionID string    `json:"transaction_id" db:"transaction_id"`
	Type          string    `json:"type" db:"type"` // DEBIT, CREDIT, BENEFIT_USAGE, REFUND
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"`
	Description   string    `json:"description" db:"description"`
	Status        string    `json:"status" db:"status"` // PENDING, COMPLETED, FAILED
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
	ClaimID       *string   `json:"claim_id,omitempty" db:"claim_id"`
	ProviderID    *string   `json:"provider_id,omitempty" db:"provider_id"`
	BlockchainTx  *string   `json:"blockchain_tx,omitempty" db:"blockchain_tx"`
}

// WalletBalance represents the current balance of a member's wallet
type WalletBalance struct {
	MemberID         string    `json:"member_id" db:"member_id"`
	TotalBalance     float64   `json:"total_balance" db:"total_balance"`
	AvailableBalance float64   `json:"available_balance" db:"available_balance"`
	ReservedBalance  float64   `json:"reserved_balance" db:"reserved_balance"`
	Currency         string    `json:"currency" db:"currency"`
	LastUpdated      time.Time `json:"last_updated" db:"last_updated"`
}

// BenefitUtilization tracks usage of specific benefits
type BenefitUtilization struct {
	ID               string    `json:"id" db:"id"`
	MemberID         string    `json:"member_id" db:"member_id"`
	BenefitType      string    `json:"benefit_type" db:"benefit_type"`
	ServiceCategory  string    `json:"service_category" db:"service_category"`
	LimitAmount      float64   `json:"limit_amount" db:"limit_amount"`
	UsedAmount       float64   `json:"used_amount" db:"used_amount"`
	RemainingAmount  float64   `json:"remaining_amount" db:"remaining_amount"`
	PeriodStart      time.Time `json:"period_start" db:"period_start"`
	PeriodEnd        time.Time `json:"period_end" db:"period_end"`
	LastUsed         *time.Time `json:"last_used,omitempty" db:"last_used"`
}

// Consent represents patient consent for data sharing and services
type Consent struct {
	ID          string    `json:"id" db:"id"`
	MemberID    string    `json:"member_id" db:"member_id"`
	Type        string    `json:"type" db:"type"` // DATA_SHARING, TELEMEDICINE, RESEARCH
	Scope       string    `json:"scope" db:"scope"`
	Status      string    `json:"status" db:"status"` // ACTIVE, REVOKED, EXPIRED
	GrantedAt   time.Time `json:"granted_at" db:"granted_at"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty" db:"expires_at"`
	RevokedAt   *time.Time `json:"revoked_at,omitempty" db:"revoked_at"`
	Purpose     string    `json:"purpose" db:"purpose"`
	DataTypes   []string  `json:"data_types"`
	Recipients  []string  `json:"recipients"`
	BlockchainHash *string `json:"blockchain_hash,omitempty" db:"blockchain_hash"`
}

// BlockchainAnchor represents data anchored to blockchain
type BlockchainAnchor struct {
	ID            string    `json:"id" db:"id"`
	RefType       string    `json:"ref_type" db:"ref_type"` // CLAIM, CONSENT, TRANSACTION
	RefID         string    `json:"ref_id" db:"ref_id"`
	Hash          string    `json:"hash" db:"hash"`
	BlockchainTx  string    `json:"blockchain_tx" db:"blockchain_tx"`
	BlockNumber   *int64    `json:"block_number,omitempty" db:"block_number"`
	Timestamp     time.Time `json:"timestamp" db:"timestamp"`
	Status        string    `json:"status" db:"status"` // PENDING, CONFIRMED, FAILED
	Signer        string    `json:"signer" db:"signer"`
}

// CostEstimate represents estimated costs for medical services
type CostEstimate struct {
	ID              string    `json:"id" db:"id"`
	MemberID        string    `json:"member_id" db:"member_id"`
	ProviderID      string    `json:"provider_id" db:"provider_id"`
	ServiceCodes    []string  `json:"service_codes"`
	EstimatedCost   float64   `json:"estimated_cost" db:"estimated_cost"`
	CoveredAmount   float64   `json:"covered_amount" db:"covered_amount"`
	PatientShare    float64   `json:"patient_share" db:"patient_share"`
	Deductible      float64   `json:"deductible" db:"deductible"`
	Copay           float64   `json:"copay" db:"copay"`
	Coinsurance     float64   `json:"coinsurance" db:"coinsurance"`
	ValidUntil      time.Time `json:"valid_until" db:"valid_until"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	RequiresPriorAuth bool    `json:"requires_prior_auth" db:"requires_prior_auth"`
}

// CreateTransactionRequest represents a request to create a new wallet transaction
type CreateTransactionRequest struct {
	Type        string  `json:"type" validate:"required,oneof=DEBIT CREDIT BENEFIT_USAGE REFUND"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description" validate:"required"`
	ClaimID     *string `json:"claim_id,omitempty"`
	ProviderID  *string `json:"provider_id,omitempty"`
}

// CreateConsentRequest represents a request to create a new consent
type CreateConsentRequest struct {
	Type       string    `json:"type" validate:"required,oneof=DATA_SHARING TELEMEDICINE RESEARCH"`
	Scope      string    `json:"scope" validate:"required"`
	Purpose    string    `json:"purpose" validate:"required"`
	DataTypes  []string  `json:"data_types" validate:"required"`
	Recipients []string  `json:"recipients" validate:"required"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
}

// CostEstimateRequest represents a request for cost estimation
type CostEstimateRequest struct {
	ProviderID   string   `json:"provider_id" validate:"required"`
	ServiceCodes []string `json:"service_codes" validate:"required"`
	ServiceDate  string   `json:"service_date" validate:"required"`
}

// BlockchainAnchorRequest represents a request to anchor data to blockchain
type BlockchainAnchorRequest struct {
	RefType string `json:"ref_type" validate:"required,oneof=CLAIM CONSENT TRANSACTION"`
	RefID   string `json:"ref_id" validate:"required"`
	Data    string `json:"data" validate:"required"`
}