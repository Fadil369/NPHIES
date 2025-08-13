package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Fadil369/NPHIES/services/eligibility-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CheckEligibility godoc
// @Summary Check member eligibility
// @Description Check eligibility and coverage for a member
// @Tags eligibility
// @Accept json
// @Produce json
// @Param request body models.EligibilityRequest true "Eligibility check request"
// @Success 200 {object} models.EligibilityResponse
// @Failure 400 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/eligibility/check [post]
func (h *Handler) CheckEligibility(c *gin.Context) {
	start := time.Now()
	h.metrics.EligibilityChecks.Inc()

	var req models.EligibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Type:    "error",
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	// Generate request ID if not provided
	if req.RequestID == "" {
		req.RequestID = uuid.New().String()
	}
	req.RequestTime = time.Now()

	// Create cache key
	cacheKey := fmt.Sprintf("eligibility:%s:%s:%s", req.MemberID, req.ProviderID, req.ServiceDate)

	// Check cache first
	var response models.EligibilityResponse
	cached, err := h.cache.Get(c.Request.Context(), cacheKey)
	if err == nil && cached != "" {
		if err := json.Unmarshal([]byte(cached), &response); err == nil {
			h.metrics.CacheHits.Inc()
			response.CacheHit = true
			response.RequestID = req.RequestID
			response.ResponseTime = time.Now()

			h.logAuditEvent(c.Request.Context(), "eligibility.check", req.RequestedBy, c.ClientIP(), map[string]interface{}{
				"request_id": req.RequestID,
				"member_id":  req.MemberID,
				"cache_hit":  true,
				"duration":   time.Since(start).Milliseconds(),
			})

			c.JSON(http.StatusOK, response)
			return
		}
	}

	h.metrics.CacheMisses.Inc()

	// Perform eligibility check
	response, err = h.performEligibilityCheck(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorf("Failed to perform eligibility check: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "ELIGIBILITY_CHECK_FAILED",
			Message: "Failed to check eligibility",
			Details: err.Error(),
		})
		return
	}

	response.CacheHit = false
	response.RequestID = req.RequestID
	response.ResponseTime = time.Now()

	// Cache the response
	responseData, _ := json.Marshal(response)
	_ = h.cache.Set(c.Request.Context(), cacheKey, string(responseData))

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "eligibility.check", req.RequestedBy, c.ClientIP(), map[string]interface{}{
		"request_id": req.RequestID,
		"member_id":  req.MemberID,
		"eligible":   response.Eligible,
		"cache_hit":  false,
		"duration":   time.Since(start).Milliseconds(),
	})

	// Check if response time exceeds SLA
	duration := time.Since(start)
	if duration.Milliseconds() > int64(h.config.Business.MaxResponseTime) {
		h.logger.Warnf("Eligibility check exceeded SLA: %dms (max: %dms)", 
			duration.Milliseconds(), h.config.Business.MaxResponseTime)
	}

	c.JSON(http.StatusOK, response)
}

// GetMemberCoverage godoc
// @Summary Get member coverage information
// @Description Retrieve detailed coverage information for a member
// @Tags eligibility
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Param effective_date query string false "Effective date for coverage lookup (YYYY-MM-DD)"
// @Success 200 {array} models.Coverage
// @Failure 404 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/eligibility/member/{id}/coverage [get]
func (h *Handler) GetMemberCoverage(c *gin.Context) {
	memberID := c.Param("id")
	effectiveDate := c.Query("effective_date")

	if effectiveDate == "" {
		effectiveDate = time.Now().Format("2006-01-02")
	}

	// Create cache key
	cacheKey := fmt.Sprintf("coverage:%s:%s", memberID, effectiveDate)

	// Check cache first
	cached, err := h.cache.Get(c.Request.Context(), cacheKey)
	if err == nil && cached != "" {
		var coverages []models.Coverage
		if err := json.Unmarshal([]byte(cached), &coverages); err == nil {
			h.metrics.CacheHits.Inc()
			c.JSON(http.StatusOK, coverages)
			return
		}
	}

	h.metrics.CacheMisses.Inc()

	// Query database for coverage
	coverages, err := h.getMemberCoverageFromDB(memberID, effectiveDate)
	if err != nil {
		h.logger.Errorf("Failed to get member coverage: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "COVERAGE_LOOKUP_FAILED",
			Message: "Failed to retrieve coverage information",
		})
		return
	}

	if len(coverages) == 0 {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Type:    "information",
			Code:    "NO_COVERAGE_FOUND",
			Message: "No active coverage found for the specified member and date",
		})
		return
	}

	// Cache the response
	responseData, _ := json.Marshal(coverages)
	_ = h.cache.Set(c.Request.Context(), cacheKey, string(responseData))

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "coverage.lookup", "", c.ClientIP(), map[string]interface{}{
		"member_id":      memberID,
		"effective_date": effectiveDate,
		"coverage_count": len(coverages),
	})

	c.JSON(http.StatusOK, coverages)
}

// VerifyCoverage godoc
// @Summary Verify coverage for specific services
// @Description Verify coverage and estimate costs for specific services
// @Tags eligibility
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Param request body models.CoverageVerificationRequest true "Coverage verification request"
// @Success 200 {object} models.CoverageVerificationResponse
// @Failure 400 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/eligibility/member/{id}/coverage/verify [post]
func (h *Handler) VerifyCoverage(c *gin.Context) {
	memberID := c.Param("id")

	var req models.CoverageVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ResponseMessage{
			Type:    "error",
			Code:    "INVALID_REQUEST",
			Message: "Invalid request format",
			Details: err.Error(),
		})
		return
	}

	req.MemberID = memberID

	// Perform coverage verification
	response, err := h.performCoverageVerification(c.Request.Context(), req)
	if err != nil {
		h.logger.Errorf("Failed to verify coverage: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "COVERAGE_VERIFICATION_FAILED",
			Message: "Failed to verify coverage",
			Details: err.Error(),
		})
		return
	}

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "coverage.verify", "", c.ClientIP(), map[string]interface{}{
		"member_id":        memberID,
		"verification_id":  response.VerificationID,
		"service_codes":    req.ServiceCodes,
		"overall_status":   response.OverallStatus,
	})

	c.JSON(http.StatusOK, response)
}

// GetMemberBenefits godoc
// @Summary Get member benefits
// @Description Retrieve detailed benefit information for a member
// @Tags eligibility
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Param service_category query string false "Service category filter"
// @Success 200 {array} models.BenefitInformation
// @Failure 404 {object} models.ResponseMessage
// @Failure 500 {object} models.ResponseMessage
// @Router /api/v1/eligibility/member/{id}/benefits [get]
func (h *Handler) GetMemberBenefits(c *gin.Context) {
	memberID := c.Param("id")
	serviceCategory := c.Query("service_category")

	// Create cache key
	cacheKey := fmt.Sprintf("benefits:%s:%s", memberID, serviceCategory)

	// Check cache first
	cached, err := h.cache.Get(c.Request.Context(), cacheKey)
	if err == nil && cached != "" {
		var benefits []models.BenefitInformation
		if err := json.Unmarshal([]byte(cached), &benefits); err == nil {
			h.metrics.CacheHits.Inc()
			c.JSON(http.StatusOK, benefits)
			return
		}
	}

	h.metrics.CacheMisses.Inc()

	// Get benefits from database/business logic
	benefits, err := h.getMemberBenefits(memberID, serviceCategory)
	if err != nil {
		h.logger.Errorf("Failed to get member benefits: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "BENEFITS_LOOKUP_FAILED",
			Message: "Failed to retrieve benefit information",
		})
		return
	}

	if len(benefits) == 0 {
		c.JSON(http.StatusNotFound, models.ResponseMessage{
			Type:    "information",
			Code:    "NO_BENEFITS_FOUND",
			Message: "No benefits found for the specified member",
		})
		return
	}

	// Cache the response
	responseData, _ := json.Marshal(benefits)
	_ = h.cache.Set(c.Request.Context(), cacheKey, string(responseData))

	// Log audit event
	h.logAuditEvent(c.Request.Context(), "benefits.lookup", "", c.ClientIP(), map[string]interface{}{
		"member_id":        memberID,
		"service_category": serviceCategory,
		"benefit_count":    len(benefits),
	})

	c.JSON(http.StatusOK, benefits)
}

// GetServiceStats godoc
// @Summary Get service statistics
// @Description Retrieve service performance and usage statistics
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} models.ServiceStats
// @Router /api/v1/admin/stats [get]
func (h *Handler) GetServiceStats(c *gin.Context) {
	stats := models.ServiceStats{
		Service: "eligibility-service",
		Version: "1.0.0",
		Uptime:  time.Since(h.startTime).String(),
		RequestStats: models.RequestStatistics{
			// These would be calculated from actual metrics in production
			TotalRequests:       1000,
			RequestsPerSecond:   10.5,
			AverageResponseTime: 120.5,
			P95ResponseTime:     250.0,
			P99ResponseTime:     400.0,
			ErrorRate:           0.02,
			SuccessRate:         0.98,
		},
		CacheStats: models.CacheStatistics{
			HitRate:       0.85,
			MissRate:      0.15,
			TotalHits:     850,
			TotalMisses:   150,
			CachedEntries: 500,
			CacheSize:     "10MB",
			EvictionCount: 25,
		},
		DatabaseStats: models.DatabaseStatistics{
			ActiveConnections: h.db.Stats().OpenConnections,
			IdleConnections:   h.db.Stats().Idle,
			TotalQueries:      1200,
			SlowQueries:       5,
			AverageQueryTime:  15.2,
		},
		Dependencies: models.DependencyStatus{
			Database: h.db.Ping() == nil,
			Redis:    h.redis.Ping(c.Request.Context()).Err() == nil,
			Kafka:    true, // Simplified check
		},
	}

	c.JSON(http.StatusOK, stats)
}

// ClearCache godoc
// @Summary Clear service cache
// @Description Clear all cached data
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/admin/cache/clear [post]
func (h *Handler) ClearCache(c *gin.Context) {
	if err := h.redis.FlushAll(c.Request.Context()).Err(); err != nil {
		h.logger.Errorf("Failed to clear cache: %v", err)
		c.JSON(http.StatusInternalServerError, models.ResponseMessage{
			Type:    "error",
			Code:    "CACHE_CLEAR_FAILED",
			Message: "Failed to clear cache",
		})
		return
	}

	h.logAuditEvent(c.Request.Context(), "cache.clear", "", c.ClientIP(), map[string]interface{}{
		"action": "clear_all",
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Cache cleared successfully",
		"status":  "success",
	})
}

// GetCacheStats godoc
// @Summary Get cache statistics
// @Description Retrieve detailed cache statistics
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} models.CacheStatistics
// @Router /api/v1/admin/cache/stats [get]
func (h *Handler) GetCacheStats(c *gin.Context) {
	info := h.redis.Info(c.Request.Context(), "memory", "stats")
	
	// Parse Redis info (simplified)
	stats := models.CacheStatistics{
		HitRate:       0.85, // Would be calculated from Redis stats
		MissRate:      0.15,
		TotalHits:     850,
		TotalMisses:   150,
		CachedEntries: 500,
		CacheSize:     "10MB",
		EvictionCount: 25,
	}

	_ = info // Use Redis info to populate actual stats

	c.JSON(http.StatusOK, stats)
}