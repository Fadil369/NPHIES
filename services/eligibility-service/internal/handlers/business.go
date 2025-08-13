package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Fadil369/NPHIES/services/eligibility-service/internal/models"
	"github.com/google/uuid"
)

// performEligibilityCheck performs the core eligibility checking logic
func (h *Handler) performEligibilityCheck(ctx context.Context, req models.EligibilityRequest) (models.EligibilityResponse, error) {
	h.metrics.DatabaseQueries.WithLabelValues("eligibility_check").Inc()

	// Get member information
	member, err := h.getMember(req.MemberID)
	if err != nil {
		return models.EligibilityResponse{}, fmt.Errorf("failed to get member: %w", err)
	}

	if member == nil {
		return models.EligibilityResponse{
			MemberID:       req.MemberID,
			Eligible:       false,
			CoverageStatus: "member_not_found",
			Messages: []models.ResponseMessage{
				{
					Type:    "error",
					Code:    "MEMBER_NOT_FOUND",
					Message: "Member not found in the system",
				},
			},
		}, nil
	}

	// Get active coverage for the member
	coverages, err := h.getMemberCoverageFromDB(req.MemberID, req.ServiceDate)
	if err != nil {
		return models.EligibilityResponse{}, fmt.Errorf("failed to get coverage: %w", err)
	}

	if len(coverages) == 0 {
		return models.EligibilityResponse{
			MemberID:       req.MemberID,
			Eligible:       false,
			CoverageStatus: "no_coverage",
			Messages: []models.ResponseMessage{
				{
					Type:    "information",
					Code:    "NO_ACTIVE_COVERAGE",
					Message: "No active coverage found for the specified date",
				},
			},
		}, nil
	}

	// Process each coverage to determine eligibility
	var benefits []models.BenefitInformation
	var limitations []models.CoverageLimitation
	var messages []models.ResponseMessage
	eligible := false
	coverageStatus := "active"

	for _, coverage := range coverages {
		// Check if coverage is active for the service date
		serviceDate, _ := time.Parse("2006-01-02", req.ServiceDate)
		if serviceDate.Before(coverage.EffectiveDate) || 
		   (coverage.ExpirationDate != nil && serviceDate.After(*coverage.ExpirationDate)) {
			continue
		}

		eligible = true

		// Calculate benefits based on coverage details
		benefitInfo := h.calculateBenefits(coverage, req.ServiceCodes, req.ProviderID)
		benefits = append(benefits, benefitInfo...)

		// Calculate limitations
		limitationInfo := h.calculateLimitations(coverage, req.ServiceCodes)
		limitations = append(limitations, limitationInfo...)

		// Add coverage-specific messages
		if coverage.Status != "active" {
			messages = append(messages, models.ResponseMessage{
				Type:    "warning",
				Code:    "COVERAGE_STATUS_WARNING",
				Message: fmt.Sprintf("Coverage status is %s", coverage.Status),
			})
		}
	}

	// Determine effective and expiration dates
	var effectiveDate, expirationDate string
	if len(coverages) > 0 {
		effectiveDate = coverages[0].EffectiveDate.Format("2006-01-02")
		if coverages[0].ExpirationDate != nil {
			expirationDate = coverages[0].ExpirationDate.Format("2006-01-02")
		}
	}

	return models.EligibilityResponse{
		MemberID:       req.MemberID,
		Eligible:       eligible,
		CoverageStatus: coverageStatus,
		EffectiveDate:  effectiveDate,
		ExpirationDate: expirationDate,
		Benefits:       benefits,
		Limitations:    limitations,
		Messages:       messages,
	}, nil
}

// performCoverageVerification performs detailed coverage verification
func (h *Handler) performCoverageVerification(ctx context.Context, req models.CoverageVerificationRequest) (models.CoverageVerificationResponse, error) {
	verificationID := uuid.New().String()

	// Get coverage information
	coverages, err := h.getMemberCoverageFromDB(req.MemberID, req.ServiceDate)
	if err != nil {
		return models.CoverageVerificationResponse{}, err
	}

	var services []models.ServiceVerification
	overallStatus := "not_covered"
	authRequired := false

	if len(coverages) > 0 {
		overallStatus = "covered"
		
		for _, serviceCode := range req.ServiceCodes {
			verification := h.verifyService(coverages[0], serviceCode, req.ProviderID, req.PlaceOfService)
			services = append(services, verification)
			
			if verification.AuthRequired {
				authRequired = true
			}
		}
	}

	return models.CoverageVerificationResponse{
		MemberID:       req.MemberID,
		VerificationID: verificationID,
		ServiceDate:    req.ServiceDate,
		Services:       services,
		OverallStatus:  overallStatus,
		AuthRequired:   authRequired,
		ValidUntil:     time.Now().Add(24 * time.Hour), // Valid for 24 hours
		Messages:       []models.ResponseMessage{},
	}, nil
}

// getMember retrieves member information from database
func (h *Handler) getMember(memberID string) (*models.Member, error) {
	query := `
		SELECT id, identifier, name, birth_date, gender, contact_info, address, status, created_at, updated_at
		FROM members 
		WHERE identifier = $1 AND status = 'active'
	`

	var member models.Member
	var nameJSON, contactJSON, addressJSON []byte

	err := h.db.QueryRow(query, memberID).Scan(
		&member.ID,
		&member.Identifier,
		&nameJSON,
		&member.BirthDate,
		&member.Gender,
		&contactJSON,
		&addressJSON,
		&member.Status,
		&member.CreatedAt,
		&member.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON fields
	json.Unmarshal(nameJSON, &member.Name)
	json.Unmarshal(contactJSON, &member.ContactInfo)
	json.Unmarshal(addressJSON, &member.Address)

	return &member, nil
}

// getMemberCoverageFromDB retrieves coverage information from database
func (h *Handler) getMemberCoverageFromDB(memberID, serviceDate string) ([]models.Coverage, error) {
	query := `
		SELECT id, member_id, payer_id, policy_number, group_number, status, type,
		       effective_date, expiration_date, benefit_details, cost_sharing,
		       network, prior_auth_rules, limitations, created_at, updated_at
		FROM coverage 
		WHERE member_id = $1 
		  AND status = 'active'
		  AND effective_date <= $2
		  AND (expiration_date IS NULL OR expiration_date >= $2)
		ORDER BY effective_date DESC
	`

	rows, err := h.db.Query(query, memberID, serviceDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coverages []models.Coverage

	for rows.Next() {
		var coverage models.Coverage
		var benefitJSON, costSharingJSON, authRulesJSON, limitationsJSON []byte

		err := rows.Scan(
			&coverage.ID,
			&coverage.MemberID,
			&coverage.PayerID,
			&coverage.PolicyNumber,
			&coverage.GroupNumber,
			&coverage.Status,
			&coverage.Type,
			&coverage.EffectiveDate,
			&coverage.ExpirationDate,
			&benefitJSON,
			&costSharingJSON,
			&coverage.Network,
			&authRulesJSON,
			&limitationsJSON,
			&coverage.CreatedAt,
			&coverage.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Parse JSON fields
		json.Unmarshal(benefitJSON, &coverage.BenefitDetails)
		json.Unmarshal(costSharingJSON, &coverage.CostSharing)
		json.Unmarshal(authRulesJSON, &coverage.PriorAuthRules)
		json.Unmarshal(limitationsJSON, &coverage.Limitations)

		coverages = append(coverages, coverage)
	}

	return coverages, nil
}

// calculateBenefits calculates benefit information based on coverage and services
func (h *Handler) calculateBenefits(coverage models.Coverage, serviceCodes []string, providerID string) []models.BenefitInformation {
	var benefits []models.BenefitInformation

	// Mock benefit calculation - in production this would involve complex business rules
	benefit := models.BenefitInformation{
		ServiceCategory:     "medical",
		InNetwork:           h.isProviderInNetwork(providerID, coverage.Network),
		CopayAmount:         25.00,
		CoinsuranceRate:     0.20,
		DeductibleAmount:    500.00,
		DeductibleMet:       false,
		RemainingDeductible: 500.00,
		OutOfPocketMax:      2000.00,
		RemainingOOPMax:     1500.00,
		PriorAuthRequired:   h.requiresPriorAuth(serviceCodes, coverage.PriorAuthRules),
		CoverageLevel:       "individual",
	}

	benefits = append(benefits, benefit)

	return benefits
}

// calculateLimitations calculates coverage limitations
func (h *Handler) calculateLimitations(coverage models.Coverage, serviceCodes []string) []models.CoverageLimitation {
	var limitations []models.CoverageLimitation

	// Mock limitation calculation
	limitation := models.CoverageLimitation{
		ServiceCategory: "medical",
		LimitationType:  "annual_maximum",
		LimitValue:      5000.00,
		UsedAmount:      1200.00,
		RemainingAmount: 3800.00,
		Period:          "annual",
		ResetDate:       "2026-01-01",
	}

	limitations = append(limitations, limitation)

	return limitations
}

// getMemberBenefits retrieves detailed benefit information
func (h *Handler) getMemberBenefits(memberID, serviceCategory string) ([]models.BenefitInformation, error) {
	// Get active coverage
	coverages, err := h.getMemberCoverageFromDB(memberID, time.Now().Format("2006-01-02"))
	if err != nil {
		return nil, err
	}

	var benefits []models.BenefitInformation

	for _, coverage := range coverages {
		// Extract benefits from coverage details
		if coverage.BenefitDetails != nil {
			// In production, this would parse the actual benefit structure
			benefit := models.BenefitInformation{
				ServiceCategory:     serviceCategory,
				InNetwork:           true,
				CopayAmount:         25.00,
				CoinsuranceRate:     0.20,
				DeductibleAmount:    500.00,
				DeductibleMet:       false,
				RemainingDeductible: 300.00,
				OutOfPocketMax:      2000.00,
				RemainingOOPMax:     1200.00,
				PriorAuthRequired:   false,
				CoverageLevel:       "individual",
			}
			benefits = append(benefits, benefit)
		}
	}

	return benefits, nil
}

// verifyService verifies coverage for a specific service
func (h *Handler) verifyService(coverage models.Coverage, serviceCode, providerID, placeOfService string) models.ServiceVerification {
	// Mock service verification - in production this would involve complex rules
	return models.ServiceVerification{
		ServiceCode:   serviceCode,
		Status:        "covered",
		CoverageLevel: 0.80, // 80% covered
		EstimatedCost: 150.00,
		PatientCost:   30.00,
		AuthRequired:  h.requiresPriorAuth([]string{serviceCode}, coverage.PriorAuthRules),
		ReasonCodes:   []string{},
	}
}

// isProviderInNetwork checks if provider is in the coverage network
func (h *Handler) isProviderInNetwork(providerID, network string) bool {
	// Mock network check - in production this would query provider network data
	return true
}

// requiresPriorAuth checks if services require prior authorization
func (h *Handler) requiresPriorAuth(serviceCodes []string, authRules map[string]interface{}) bool {
	// Mock prior auth check - in production this would evaluate complex rules
	for _, code := range serviceCodes {
		// Check if code requires prior auth
		if code == "99245" || code == "99244" { // Example codes that require auth
			return true
		}
	}
	return false
}