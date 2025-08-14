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

	"github.com/Fadil369/NPHIES/services/terminology-service/internal/config"
	"github.com/Fadil369/NPHIES/services/terminology-service/internal/handlers"
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
	h := handlers.NewHandler(logger)

	// Setup router
	router := setupRouter(cfg, h, logger)

	// Setup server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server
	go func() {
		logger.WithField("port", cfg.Server.Port).Info("Starting terminology service")
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
		// Code Systems
		v1.GET("/codesystems", h.GetCodeSystems)
		v1.GET("/codesystems/:system/codes/:code", h.LookupCode)
		v1.POST("/codesystems/:system/validate", h.ValidateCode)
		v1.GET("/codesystems/:system/search", h.SearchCodes)

		// Terminology mapping
		v1.POST("/terminology/map", h.MapCodes)
		v1.GET("/terminology/concepts/:concept", h.GetConcept)
	}

	return router
}