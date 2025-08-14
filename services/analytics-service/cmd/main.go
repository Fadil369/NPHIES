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

	"github.com/Fadil369/NPHIES/services/analytics-service/internal/config"
	"github.com/Fadil369/NPHIES/services/analytics-service/internal/handlers"
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
		logger.WithField("port", cfg.Server.Port).Info("Starting analytics service")
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
		// Predictive Analytics endpoints
		v1.POST("/analytics/risk-stratification", h.RiskStratification)
		v1.POST("/analytics/cost-forecast", h.CostForecast)
		v1.POST("/analytics/fraud-detection", h.FraudDetection)
		v1.GET("/analytics/population-health", h.PopulationHealth)

		// IoT Integration endpoints
		v1.POST("/iot/data-ingestion", h.IngestIoTData)
		v1.GET("/iot/devices/:memberId", h.GetMemberDevices)
		v1.GET("/iot/metrics/:deviceId", h.GetDeviceMetrics)

		// Advanced Privacy endpoints
		v1.POST("/privacy/differential-export", h.DifferentialPrivacyExport)
		v1.POST("/privacy/anonymize", h.AnonymizeData)
		v1.GET("/privacy/compliance-check", h.ComplianceCheck)

		// Real-time Analytics
		v1.GET("/analytics/realtime/claims", h.RealtimeClaimsAnalytics)
		v1.GET("/analytics/realtime/utilization", h.RealtimeUtilization)
		v1.GET("/analytics/trends/:period", h.AnalyticsTrends)

		// ML Model endpoints
		v1.POST("/ml/models/:modelId/predict", h.PredictWithModel)
		v1.GET("/ml/models", h.ListModels)
		v1.POST("/ml/models/:modelId/retrain", h.RetrainModel)
		v1.GET("/ml/models/:modelId/performance", h.GetModelPerformance)
	}

	return router
}