package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Fadil369/NPHIES/services/eligibility-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// SearchCoverage godoc
// @Summary Search coverage records
// @Description Search for coverage records with various filters
// @Tags coverage
// @Accept json
// @Produce json
// @Param member_id query string false "Member ID"
// @Param payer_id query string false "Payer ID"
// @Param status query string false "Coverage status"
// @Param effective_date query string false "Effective date (YYYY-MM-DD)"
// @Param _count query int false "Number of results to return" default(20)
// @Param _offset query int false "Offset for pagination" default(0)
// @Success 200 {array} models.Coverage
// @Failure 400 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/coverage [get]
func (h *Handler) SearchCoverage(c *gin.Context) {
	// Parse query parameters
	memberID := c.Query("member_id")
	payerID := c.Query("payer_id")
	status := c.Query("status")
	effectiveDate := c.Query("effective_date")

	// Parse pagination parameters
	count := 20
	if countStr := c.Query("_count"); countStr != "" {
		if parsedCount, err := strconv.Atoi(countStr); err == nil && parsedCount > 0 && parsedCount <= 100 {
			count = parsedCount
		}
	}

	offset := 0
	if offsetStr := c.Query("_offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	// Build query
	query := `
		SELECT id, member_id, payer_id, policy_number, group_number, status, type,
		       effective_date, expiration_date, benefit_details, cost_sharing,
		       network, prior_auth_rules, limitations, created_at, updated_at
		FROM coverage 
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 1

	if memberID != "" {
		query += " AND member_id = $" + strconv.Itoa(argIndex)
		args = append(args, memberID)
		argIndex++
	}

	if payerID != "" {
		query += " AND payer_id = $" + strconv.Itoa(argIndex)
		args = append(args, payerID)
		argIndex++
	}

	if status != "" {
		query += " AND status = $" + strconv.Itoa(argIndex)
		args = append(args, status)
		argIndex++
	}

	if effectiveDate != "" {
		query += " AND effective_date <= $" + strconv.Itoa(argIndex)
		args = append(args, effectiveDate)
		argIndex++
	}

	query += " ORDER BY created_at DESC LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, count, offset)

	rows, err := h.db.Query(query, args...)
	if err != nil {
		h.logger.Errorf("Failed to search coverage: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "COVERAGE_SEARCH_FAILED",
			Message: "Failed to search coverage records",
		})
		return
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
			h.logger.Errorf("Failed to scan coverage row: %v", err)
			continue
		}

		// Parse JSON fields
		json.Unmarshal(benefitJSON, &coverage.BenefitDetails)
		json.Unmarshal(costSharingJSON, &coverage.CostSharing)
		json.Unmarshal(authRulesJSON, &coverage.PriorAuthRules)
		json.Unmarshal(limitationsJSON, &coverage.Limitations)

		coverages = append(coverages, coverage)
	}

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "coverage.search", "", c.ClientIP(), map[string]interface{}{
		"filters": map[string]string{
			"member_id":      memberID,
			"payer_id":       payerID,
			"status":         status,
			"effective_date": effectiveDate,
		},
		"result_count": len(coverages),
		"offset":       offset,
		"limit":        count,
	})

	c.JSON(http.StatusOK, coverages)
}

// CreateCoverage godoc
// @Summary Create a new coverage record
// @Description Create a new coverage record for a member
// @Tags coverage
// @Accept json
// @Produce json
// @Param coverage body models.Coverage true "Coverage record"
// @Success 201 {object} models.Coverage
// @Failure 400 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/coverage [post]
func (h *Handler) CreateCoverage(c *gin.Context) {
	var coverage models.Coverage
	if err := c.ShouldBindJSON(&coverage); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Type:    "error",
			Code:    "INVALID_REQUEST",
			Message: "Invalid coverage data",
			Details: err.Error(),
		})
		return
	}

	// Validate required fields
	if coverage.MemberID == "" || coverage.PayerID == "" {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Type:    "error",
			Code:    "MISSING_REQUIRED_FIELDS",
			Message: "Member ID and Payer ID are required",
		})
		return
	}

	// Generate ID and set timestamps
	coverage.ID = uuid.New().String()
	coverage.CreatedAt = time.Now()
	coverage.UpdatedAt = time.Now()

	if coverage.Status == "" {
		coverage.Status = "active"
	}

	// Serialize JSON fields
	benefitJSON, _ := json.Marshal(coverage.BenefitDetails)
	costSharingJSON, _ := json.Marshal(coverage.CostSharing)
	authRulesJSON, _ := json.Marshal(coverage.PriorAuthRules)
	limitationsJSON, _ := json.Marshal(coverage.Limitations)

	// Insert into database
	query := `
		INSERT INTO coverage (
			id, member_id, payer_id, policy_number, group_number, status, type,
			effective_date, expiration_date, benefit_details, cost_sharing,
			network, prior_auth_rules, limitations, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`

	_, err := h.db.Exec(query,
		coverage.ID,
		coverage.MemberID,
		coverage.PayerID,
		coverage.PolicyNumber,
		coverage.GroupNumber,
		coverage.Status,
		coverage.Type,
		coverage.EffectiveDate,
		coverage.ExpirationDate,
		benefitJSON,
		costSharingJSON,
		coverage.Network,
		authRulesJSON,
		limitationsJSON,
		coverage.CreatedAt,
		coverage.UpdatedAt,
	)

	if err != nil {
		h.logger.Errorf("Failed to create coverage: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "COVERAGE_CREATE_FAILED",
			Message: "Failed to create coverage record",
		})
		return
	}

	// Clear cache for this member
	cachePattern := fmt.Sprintf("*:%s:*", coverage.MemberID)
	_ = h.cache.DeletePattern(c.Request.Context(), cachePattern)

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "coverage.create", "", c.ClientIP(), map[string]interface{}{
		"coverage_id": coverage.ID,
		"member_id":   coverage.MemberID,
		"payer_id":    coverage.PayerID,
		"status":      coverage.Status,
	})

	c.Header("Location", "/api/v1/coverage/"+coverage.ID)
	c.JSON(http.StatusCreated, coverage)
}

// GetCoverage godoc
// @Summary Get coverage by ID
// @Description Retrieve a specific coverage record by ID
// @Tags coverage
// @Accept json
// @Produce json
// @Param id path string true "Coverage ID"
// @Success 200 {object} models.Coverage
// @Failure 404 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/coverage/{id} [get]
func (h *Handler) GetCoverage(c *gin.Context) {
	coverageID := c.Param("id")

	query := `
		SELECT id, member_id, payer_id, policy_number, group_number, status, type,
		       effective_date, expiration_date, benefit_details, cost_sharing,
		       network, prior_auth_rules, limitations, created_at, updated_at
		FROM coverage 
		WHERE id = $1
	`

	var coverage models.Coverage
	var benefitJSON, costSharingJSON, authRulesJSON, limitationsJSON []byte

	err := h.db.QueryRow(query, coverageID).Scan(
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
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, models.ResponseMessage{
				Type:    "error",
				Code:    "COVERAGE_NOT_FOUND",
				Message: "Coverage record not found",
			})
		} else {
			h.logger.Errorf("Failed to get coverage: %v", err)
			c.JSON(http.StatusInternalServerError, models.ResponseMessage{
				Type:    "error",
				Code:    "COVERAGE_RETRIEVAL_FAILED",
				Message: "Failed to retrieve coverage record",
			})
		}
		return
	}

	// Parse JSON fields
	json.Unmarshal(benefitJSON, &coverage.BenefitDetails)
	json.Unmarshal(costSharingJSON, &coverage.CostSharing)
	json.Unmarshal(authRulesJSON, &coverage.PriorAuthRules)
	json.Unmarshal(limitationsJSON, &coverage.Limitations)

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "coverage.read", "", c.ClientIP(), map[string]interface{}{
		"coverage_id": coverageID,
		"member_id":   coverage.MemberID,
	})

	c.JSON(http.StatusOK, coverage)
}

// UpdateCoverage godoc
// @Summary Update coverage record
// @Description Update an existing coverage record
// @Tags coverage
// @Accept json
// @Produce json
// @Param id path string true "Coverage ID"
// @Param coverage body models.Coverage true "Updated coverage record"
// @Success 200 {object} models.Coverage
// @Failure 400 {object} models.ResponseMessage
// @Failure 404 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/coverage/{id} [put]
func (h *Handler) UpdateCoverage(c *gin.Context) {
	coverageID := c.Param("id")

	var coverage models.Coverage
	if err := c.ShouldBindJSON(&coverage); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Type:    "error",
			Code:    "INVALID_REQUEST",
			Message: "Invalid coverage data",
			Details: err.Error(),
		})
		return
	}

	// Ensure ID matches
	coverage.ID = coverageID
	coverage.UpdatedAt = time.Now()

	// Serialize JSON fields
	benefitJSON, _ := json.Marshal(coverage.BenefitDetails)
	costSharingJSON, _ := json.Marshal(coverage.CostSharing)
	authRulesJSON, _ := json.Marshal(coverage.PriorAuthRules)
	limitationsJSON, _ := json.Marshal(coverage.Limitations)

	// Update database
	query := `
		UPDATE coverage SET
			member_id = $2, payer_id = $3, policy_number = $4, group_number = $5,
			status = $6, type = $7, effective_date = $8, expiration_date = $9,
			benefit_details = $10, cost_sharing = $11, network = $12,
			prior_auth_rules = $13, limitations = $14, updated_at = $15
		WHERE id = $1
	`

	result, err := h.db.Exec(query,
		coverage.ID,
		coverage.MemberID,
		coverage.PayerID,
		coverage.PolicyNumber,
		coverage.GroupNumber,
		coverage.Status,
		coverage.Type,
		coverage.EffectiveDate,
		coverage.ExpirationDate,
		benefitJSON,
		costSharingJSON,
		coverage.Network,
		authRulesJSON,
		limitationsJSON,
		coverage.UpdatedAt,
	)

	if err != nil {
		h.logger.Errorf("Failed to update coverage: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "COVERAGE_UPDATE_FAILED",
			Message: "Failed to update coverage record",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Type:    "error",
			Code:    "COVERAGE_NOT_FOUND",
			Message: "Coverage record not found",
		})
		return
	}

	// Clear cache for this member
	cachePattern := fmt.Sprintf("*:%s:*", coverage.MemberID)
	_ = h.cache.DeletePattern(c.Request.Context(), cachePattern)

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "coverage.update", "", c.ClientIP(), map[string]interface{}{
		"coverage_id": coverageID,
		"member_id":   coverage.MemberID,
	})

	c.JSON(http.StatusOK, coverage)
}

// DeleteCoverage godoc
// @Summary Delete coverage record
// @Description Delete a coverage record (soft delete by changing status)
// @Tags coverage
// @Accept json
// @Produce json
// @Param id path string true "Coverage ID"
// @Success 204
// @Failure 404 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/coverage/{id} [delete]
func (h *Handler) DeleteCoverage(c *gin.Context) {
	coverageID := c.Param("id")

	// Soft delete by updating status
	query := `
		UPDATE coverage SET 
			status = 'deleted', 
			updated_at = $2 
		WHERE id = $1 AND status != 'deleted'
	`

	result, err := h.db.Exec(query, coverageID, time.Now())
	if err != nil {
		h.logger.Errorf("Failed to delete coverage: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "COVERAGE_DELETE_FAILED",
			Message: "Failed to delete coverage record",
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Type:    "error",
			Code:    "COVERAGE_NOT_FOUND",
			Message: "Coverage record not found",
		})
		return
	}

	// Clear cache (we don't know the member ID, so clear by coverage ID pattern)
	cachePattern := fmt.Sprintf("*:%s", coverageID)
	_ = h.cache.DeletePattern(c.Request.Context(), cachePattern)

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "coverage.delete", "", c.ClientIP(), map[string]interface{}{
		"coverage_id": coverageID,
	})

	c.Status(http.StatusNoContent)
}