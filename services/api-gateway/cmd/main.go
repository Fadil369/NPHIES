package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fadil369/NPHIES/services/api-gateway/internal/config"
	"github.com/Fadil369/NPHIES/services/api-gateway/internal/handlers"
	"github.com/Fadil369/NPHIES/services/api-gateway/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title NPHIES API Gateway
// @version 1.0
// @description Unified Digital Healthcare Insurance Platform API Gateway
// @termsOfService https://github.com/Fadil369/NPHIES

// @contact.name NPHIES Support
// @contact.url https://github.com/Fadil369/NPHIES
// @contact.email support@nphies.sa

// @license.name MIT
// @license.url https://github.com/Fadil369/NPHIES/blob/main/LICENSE

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.oauth2.application OAuth2Application
// @tokenUrl https://auth.nphies.sa/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access
// @scope.admin Grants admin access

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
		logger.SetLevel(logrus.InfoLevel)
	}

	// Initialize handlers
	h, err := handlers.NewHandler(cfg, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize handlers: %v", err)
	}
	defer h.Close()

	// Setup router
	router := setupRouter(cfg, h, logger)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.Infof("Starting NPHIES API Gateway on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

func setupRouter(cfg *config.Config, h *handlers.Handler, logger *logrus.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.RecoveryMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.SecurityHeadersMiddleware())
	router.Use(middleware.RateLimitMiddleware(cfg.RateLimit))
	router.Use(middleware.MetricsMiddleware())

	// Health checks
	router.GET("/health", h.HealthCheck)
	router.GET("/ready", h.ReadinessCheck)
	router.GET("/metrics", h.MetricsHandler)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Authentication
		auth := v1.Group("/auth")
		{
			auth.POST("/token", h.GetToken)
			auth.POST("/refresh", h.RefreshToken)
		}

		// FHIR Resources - protected endpoints
		fhir := v1.Group("/fhir").Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// Patient endpoints
			patients := fhir.Group("/Patient")
			{
				patients.GET("", h.SearchPatients)
				patients.POST("", h.CreatePatient)
				patients.GET("/:id", h.GetPatient)
				patients.PUT("/:id", h.UpdatePatient)
				patients.DELETE("/:id", h.DeletePatient)
			}

			// Coverage endpoints
			coverage := fhir.Group("/Coverage")
			{
				coverage.GET("", h.SearchCoverage)
				coverage.POST("", h.CreateCoverage)
				coverage.GET("/:id", h.GetCoverage)
				coverage.PUT("/:id", h.UpdateCoverage)
				coverage.DELETE("/:id", h.DeleteCoverage)
			}

			// Claim endpoints
			claims := fhir.Group("/Claim")
			{
				claims.GET("", h.SearchClaims)
				claims.POST("", h.CreateClaim)
				claims.GET("/:id", h.GetClaim)
				claims.PUT("/:id", h.UpdateClaim)
				claims.DELETE("/:id", h.DeleteClaim)
			}

			// ClaimResponse endpoints
			claimResponses := fhir.Group("/ClaimResponse")
			{
				claimResponses.GET("", h.SearchClaimResponses)
				claimResponses.GET("/:id", h.GetClaimResponse)
			}

			// Prior Authorization endpoints
			priorAuth := fhir.Group("/CoverageEligibilityRequest")
			{
				priorAuth.GET("", h.SearchPriorAuthorizations)
				priorAuth.POST("", h.CreatePriorAuthorization)
				priorAuth.GET("/:id", h.GetPriorAuthorization)
				priorAuth.PUT("/:id", h.UpdatePriorAuthorization)
			}
		}

		// Eligibility Service Proxy
		eligibility := v1.Group("/eligibility").Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			eligibility.POST("/check", h.CheckEligibility)
			eligibility.GET("/member/:id/coverage", h.GetMemberCoverage)
		}

		// Claims Service Proxy
		claimsProxy := v1.Group("/claims").Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			claimsProxy.POST("/submit", h.SubmitClaim)
			claimsProxy.GET("/:id/status", h.GetClaimStatus)
			claimsProxy.POST("/:id/reprocess", h.ReprocessClaim)
		}

		// Terminology Service Proxy
		terminology := v1.Group("/terminology").Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			terminology.GET("/codesystems", h.GetCodeSystems)
			terminology.GET("/codesystems/:system/codes/:code", h.LookupCode)
			terminology.POST("/codesystems/:system/validate", h.ValidateCode)
		}

		// Administrative endpoints
		admin := v1.Group("/admin").Use(middleware.AuthMiddleware(cfg.JWT.Secret), middleware.AdminMiddleware())
		{
			admin.GET("/stats", h.GetSystemStats)
			admin.GET("/audit", h.GetAuditLogs)
			admin.POST("/cache/clear", h.ClearCache)
		}
	}

	// Swagger documentation
	// Note: In production, this should be behind authentication
	// router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}