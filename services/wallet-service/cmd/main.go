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

	"github.com/Fadil369/NPHIES/services/wallet-service/internal/config"
	"github.com/Fadil369/NPHIES/services/wallet-service/internal/handlers"
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
		logger.WithField("port", cfg.Server.Port).Info("Starting wallet service")
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
		// Digital Wallet endpoints
		v1.GET("/wallet/:memberId", h.GetWallet)
		v1.POST("/wallet/:memberId/transactions", h.CreateTransaction)
		v1.GET("/wallet/:memberId/transactions", h.GetTransactions)
		v1.GET("/wallet/:memberId/balance", h.GetBalance)

		// Blockchain anchoring endpoints
		v1.POST("/blockchain/anchor", h.AnchorToBlockchain)
		v1.GET("/blockchain/verify/:hash", h.VerifyHash)
		v1.GET("/blockchain/transaction/:txId", h.GetBlockchainTransaction)

		// Consent management
		v1.POST("/consent", h.CreateConsent)
		v1.GET("/consent/:memberId", h.GetConsents)
		v1.PUT("/consent/:consentId", h.UpdateConsent)
		v1.DELETE("/consent/:consentId", h.RevokeConsent)

		// Cost estimation
		v1.POST("/estimate/cost", h.EstimateCost)
		v1.GET("/estimate/provider/:providerId/services", h.GetProviderServices)

		// Benefits tracking
		v1.GET("/benefits/:memberId/remaining", h.GetRemainingBenefits)
		v1.POST("/benefits/:memberId/deduct", h.DeductBenefits)
		v1.GET("/benefits/:memberId/utilization", h.GetBenefitUtilization)
	}

	return router
}