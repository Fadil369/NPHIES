package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Fadil369/NPHIES/services/automation-service/internal/config"
	"github.com/Fadil369/NPHIES/services/automation-service/internal/handlers"
)

func main() {
	// Initialize logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.WithError(err).Fatal("Failed to load configuration")
	}

	// Initialize handlers
	h, err := handlers.NewHandler(logger, cfg)
	if err != nil {
		logger.WithError(err).Fatal("Failed to initialize handlers")
	}

	// Setup router
	router := setupRouter(cfg, h, logger)

	// Setup server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server
	go func() {
		logger.WithField("port", cfg.Server.Port).Info("Starting automation service")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Server exited")
}

func setupRouter(cfg *config.Config, h *handlers.Handler, logger *logrus.Logger) *gin.Engine {
	if cfg.Server.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Health endpoints
	router.GET("/health", h.Health)
	router.GET("/ready", h.Ready)
	router.GET("/metrics", h.Metrics)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Advanced AI/ML Automation
		v1.POST("/automation/ai/auto-adjudicate", h.AutoAdjudicate)
		v1.POST("/automation/ai/auto-authorize", h.AutoAuthorize)
		v1.POST("/automation/ai/intelligent-routing", h.IntelligentRouting)
		v1.POST("/automation/ai/predictive-intervention", h.PredictiveIntervention)

		// Real-time Decision Engines
		v1.POST("/decision/realtime/eligibility", h.RealtimeEligibilityDecision)
		v1.POST("/decision/realtime/pricing", h.RealtimePricingDecision)
		v1.POST("/decision/realtime/authorization", h.RealtimeAuthorizationDecision)
		v1.POST("/decision/realtime/fraud-alert", h.RealtimeFraudAlert)

		// Complete Workflow Automation
		v1.POST("/workflow/end-to-end/claims-processing", h.EndToEndClaimsProcessing)
		v1.POST("/workflow/end-to-end/member-onboarding", h.EndToEndMemberOnboarding)
		v1.POST("/workflow/end-to-end/provider-integration", h.EndToEndProviderIntegration)
		v1.POST("/workflow/end-to-end/compliance-check", h.EndToEndComplianceCheck)

		// Advanced Security & Compliance
		v1.POST("/security/adaptive-auth", h.AdaptiveAuthentication)
		v1.POST("/security/threat-detection", h.ThreatDetection)
		v1.POST("/security/compliance-automation", h.ComplianceAutomation)
		v1.GET("/security/risk-assessment", h.RiskAssessment)

		// Population Health Analytics (Complete)
		v1.GET("/analytics/population/health-trends", h.PopulationHealthTrends)
		v1.GET("/analytics/population/disease-burden", h.DiseaseBurdenAnalysis)
		v1.GET("/analytics/population/cost-drivers", h.CostDriverAnalysis)
		v1.GET("/analytics/population/quality-metrics", h.QualityMetrics)
		v1.POST("/analytics/population/intervention-opportunities", h.InterventionOpportunities)

		// Advanced AI Model Management
		v1.POST("/ai/models/deploy", h.DeployModel)
		v1.POST("/ai/models/a-b-test", h.ABTestModels)
		v1.POST("/ai/models/auto-retrain", h.AutoRetrainModels)
		v1.GET("/ai/models/performance-monitoring", h.ModelPerformanceMonitoring)
		v1.POST("/ai/models/explainability", h.ModelExplainability)

		// Intelligent Automation Orchestration
		v1.POST("/orchestration/create-workflow", h.CreateAutomatedWorkflow)
		v1.GET("/orchestration/workflows", h.ListWorkflows)
		v1.POST("/orchestration/workflows/:id/execute", h.ExecuteWorkflow)
		v1.GET("/orchestration/workflows/:id/status", h.GetWorkflowStatus)
	}

	return router
}