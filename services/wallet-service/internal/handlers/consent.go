package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/Fadil369/NPHIES/services/wallet-service/internal/blockchain"
	"github.com/Fadil369/NPHIES/services/wallet-service/internal/models"
)

// CreateConsent creates a new consent record
func (h *Handler) CreateConsent(c *gin.Context) {
	var request models.CreateConsentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	memberID := c.GetHeader("X-Member-ID")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Member ID required in header"})
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"member_id": memberID,
		"type":      request.Type,
		"scope":     request.Scope,
	}).Info("Creating consent record")

	// Create consent
	consent := models.Consent{
		ID:        uuid.New().String(),
		MemberID:  memberID,
		Type:      request.Type,
		Scope:     request.Scope,
		Status:    "ACTIVE",
		GrantedAt: time.Now(),
		ExpiresAt: request.ExpiresAt,
		Purpose:   request.Purpose,
		DataTypes: request.DataTypes,
		Recipients: request.Recipients,
	}

	// Anchor consent to blockchain
	if blockchainHash, err := h.anchorConsentToBlockchain(consent); err == nil {
		consent.BlockchainHash = &blockchainHash
	}

	c.JSON(http.StatusCreated, consent)
}

// GetConsents retrieves consent records for a member
func (h *Handler) GetConsents(c *gin.Context) {
	memberID := c.Param("memberId")

	h.logger.WithField("member_id", memberID).Info("Retrieving consent records")

	// Mock consent data
	consents := []gin.H{
		{
			"id":          "consent-001",
			"member_id":   memberID,
			"type":        "DATA_SHARING",
			"scope":       "CLAIMS_DATA",
			"status":      "ACTIVE",
			"granted_at":  time.Now().Add(-30 * 24 * time.Hour),
			"expires_at":  time.Now().Add(365 * 24 * time.Hour),
			"purpose":     "Insurance processing and analytics",
			"data_types":  []string{"CLAIMS", "ELIGIBILITY", "DEMOGRAPHICS"},
			"recipients":  []string{"PRIMARY_INSURER", "ANALYTICS_PLATFORM"},
			"blockchain_hash": "0xabcdef1234567890",
		},
		{
			"id":         "consent-002",
			"member_id":  memberID,
			"type":       "TELEMEDICINE",
			"scope":      "CONSULTATION_DATA",
			"status":     "ACTIVE",
			"granted_at": time.Now().Add(-7 * 24 * time.Hour),
			"purpose":    "Remote healthcare consultation",
			"data_types": []string{"VITAL_SIGNS", "MEDICAL_HISTORY"},
			"recipients": []string{"TELEMEDICINE_PROVIDER"},
			"blockchain_hash": "0x1234567890abcdef",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"consents": consents,
		"total":    len(consents),
	})
}

// UpdateConsent updates an existing consent record
func (h *Handler) UpdateConsent(c *gin.Context) {
	consentID := c.Param("consentId")

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"consent_id": consentID,
		"updates":    updates,
	}).Info("Updating consent record")

	// Mock updated consent
	updatedConsent := gin.H{
		"id":        consentID,
		"status":    "ACTIVE",
		"updated_at": time.Now(),
		"updates":   updates,
	}

	c.JSON(http.StatusOK, updatedConsent)
}

// RevokeConsent revokes a consent record
func (h *Handler) RevokeConsent(c *gin.Context) {
	consentID := c.Param("consentId")

	h.logger.WithField("consent_id", consentID).Info("Revoking consent record")

	// Create revocation record for blockchain
	revocation := gin.H{
		"consent_id": consentID,
		"revoked_at": time.Now(),
		"status":     "REVOKED",
	}

	// Anchor revocation to blockchain
	if blockchainTx, err := h.blockchainClient.SubmitTransaction("CONSENT_REVOCATION", consentID, "REVOKED"); err == nil {
		revocation["blockchain_tx"] = blockchainTx.TransactionID
	}

	c.JSON(http.StatusOK, revocation)
}

// EstimateCost estimates costs for medical services
func (h *Handler) EstimateCost(c *gin.Context) {
	var request models.CostEstimateRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	memberID := c.GetHeader("X-Member-ID")
	if memberID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Member ID required in header"})
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"member_id":     memberID,
		"provider_id":   request.ProviderID,
		"service_codes": request.ServiceCodes,
	}).Info("Estimating cost for services")

	// Mock cost estimation logic
	estimatedCost := 750.00
	coveredAmount := 600.00
	patientShare := 150.00

	estimate := models.CostEstimate{
		ID:              uuid.New().String(),
		MemberID:        memberID,
		ProviderID:      request.ProviderID,
		ServiceCodes:    request.ServiceCodes,
		EstimatedCost:   estimatedCost,
		CoveredAmount:   coveredAmount,
		PatientShare:    patientShare,
		Deductible:      50.00,
		Copay:           100.00,
		Coinsurance:     0.00,
		ValidUntil:      time.Now().Add(30 * 24 * time.Hour),
		CreatedAt:       time.Now(),
		RequiresPriorAuth: len(request.ServiceCodes) > 2,
	}

	c.JSON(http.StatusOK, estimate)
}

// GetProviderServices retrieves available services for a provider
func (h *Handler) GetProviderServices(c *gin.Context) {
	providerID := c.Param("providerId")

	h.logger.WithField("provider_id", providerID).Info("Retrieving provider services")

	// Mock provider services
	services := []gin.H{
		{
			"code":        "99213",
			"description": "Office visit - established patient",
			"category":    "CONSULTATION",
			"base_cost":   200.00,
			"duration":    30,
			"requires_prior_auth": false,
		},
		{
			"code":        "73060",
			"description": "Knee X-ray, 2 views",
			"category":    "IMAGING",
			"base_cost":   150.00,
			"duration":    15,
			"requires_prior_auth": false,
		},
		{
			"code":        "29881",
			"description": "Arthroscopy, knee, surgical",
			"category":    "SURGERY",
			"base_cost":   3500.00,
			"duration":    120,
			"requires_prior_auth": true,
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"provider_id": providerID,
		"services":    services,
		"total":       len(services),
	})
}

// GetRemainingBenefits retrieves remaining benefits for a member
func (h *Handler) GetRemainingBenefits(c *gin.Context) {
	memberID := c.Param("memberId")

	h.logger.WithField("member_id", memberID).Info("Retrieving remaining benefits")

	// Mock benefit utilization data
	benefits := []models.BenefitUtilization{
		{
			ID:              uuid.New().String(),
			MemberID:        memberID,
			BenefitType:     "MEDICAL",
			ServiceCategory: "GENERAL",
			LimitAmount:     5000.00,
			UsedAmount:      1200.00,
			RemainingAmount: 3800.00,
			PeriodStart:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			PeriodEnd:       time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			ID:              uuid.New().String(),
			MemberID:        memberID,
			BenefitType:     "DENTAL",
			ServiceCategory: "PREVENTIVE",
			LimitAmount:     1000.00,
			UsedAmount:      200.00,
			RemainingAmount: 800.00,
			PeriodStart:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			PeriodEnd:       time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"member_id": memberID,
		"benefits":  benefits,
		"total":     len(benefits),
	})
}

// DeductBenefits deducts benefits for a service
func (h *Handler) DeductBenefits(c *gin.Context) {
	memberID := c.Param("memberId")

	var request struct {
		BenefitType     string  `json:"benefit_type" validate:"required"`
		ServiceCategory string  `json:"service_category" validate:"required"`
		Amount          float64 `json:"amount" validate:"required,gt=0"`
		Description     string  `json:"description" validate:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	h.logger.WithFields(map[string]interface{}{
		"member_id":        memberID,
		"benefit_type":     request.BenefitType,
		"service_category": request.ServiceCategory,
		"amount":           request.Amount,
	}).Info("Deducting benefits")

	// Mock benefit deduction
	deduction := gin.H{
		"member_id":        memberID,
		"benefit_type":     request.BenefitType,
		"service_category": request.ServiceCategory,
		"deducted_amount":  request.Amount,
		"remaining_amount": 3550.00, // Mock remaining after deduction
		"description":      request.Description,
		"timestamp":        time.Now(),
		"status":           "COMPLETED",
	}

	c.JSON(http.StatusOK, deduction)
}

// GetBenefitUtilization retrieves benefit utilization history
func (h *Handler) GetBenefitUtilization(c *gin.Context) {
	memberID := c.Param("memberId")

	h.logger.WithField("member_id", memberID).Info("Retrieving benefit utilization history")

	// Mock utilization data
	utilization := []gin.H{
		{
			"date":             "2024-01-15",
			"benefit_type":     "MEDICAL",
			"service_category": "CONSULTATION",
			"amount_used":      250.00,
			"description":      "Doctor visit",
			"provider_id":      "PRV-001",
			"claim_id":         "CLM-001",
		},
		{
			"date":             "2024-01-20",
			"benefit_type":     "DENTAL",
			"service_category": "PREVENTIVE",
			"amount_used":      200.00,
			"description":      "Dental cleaning",
			"provider_id":      "PRV-002",
			"claim_id":         "CLM-002",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"member_id":   memberID,
		"utilization": utilization,
		"total":       len(utilization),
	})
}

// anchorConsentToBlockchain anchors a consent record to blockchain
func (h *Handler) anchorConsentToBlockchain(consent models.Consent) (string, error) {
	data, err := blockchain.CreateDataHash(consent)
	if err != nil {
		return "", err
	}

	response, err := h.blockchainClient.SubmitTransaction("CONSENT", consent.ID, data)
	if err != nil {
		return "", err
	}

	return response.Hash, nil
}