package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fadil369/NPHIES/services/eligibility-service/internal/config"
	"github.com/Fadil369/NPHIES/services/eligibility-service/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @title NPHIES Eligibility Service
// @version 1.0
// @description Real-time coverage validation and eligibility checking service
// @termsOfService https://github.com/Fadil369/NPHIES

// @contact.name NPHIES Support
// @contact.url https://github.com/Fadil369/NPHIES
// @contact.email support@nphies.sa

// @license.name MIT
// @license.url https://github.com/Fadil369/NPHIES/blob/main/LICENSE

// @host localhost:8090
// @BasePath /api/v1

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
	router := setupRouter(h, logger)

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
		logger.Infof("Starting NPHIES Eligibility Service on port %s", cfg.Port)
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

func setupRouter(h *handlers.Handler, logger *logrus.Logger) *gin.Engine {
	router := gin.New()

	// Middleware
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"method":     param.Method,
			"path":       param.Path,
			"status":     param.StatusCode,
			"latency":    param.Latency,
			"client_ip":  param.ClientIP,
			"user_agent": param.Request.UserAgent(),
		}).Info("HTTP Request")
		return ""
	}))
	router.Use(gin.Recovery())

	// Health checks
	router.GET("/health", h.HealthCheck)
	router.GET("/ready", h.ReadinessCheck)
	router.GET("/metrics", h.MetricsHandler)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Eligibility endpoints
		eligibility := v1.Group("/eligibility")
		{
			eligibility.POST("/check", h.CheckEligibility)
			eligibility.GET("/member/:id/coverage", h.GetMemberCoverage)
			eligibility.POST("/member/:id/coverage/verify", h.VerifyCoverage)
			eligibility.GET("/member/:id/benefits", h.GetMemberBenefits)
		}

		// Coverage endpoints
		coverage := v1.Group("/coverage")
		{
			coverage.GET("", h.SearchCoverage)
			coverage.POST("", h.CreateCoverage)
			coverage.GET("/:id", h.GetCoverage)
			coverage.PUT("/:id", h.UpdateCoverage)
			coverage.DELETE("/:id", h.DeleteCoverage)
		}

		// Administrative endpoints
		admin := v1.Group("/admin")
		{
			admin.GET("/stats", h.GetServiceStats)
			admin.POST("/cache/clear", h.ClearCache)
			admin.GET("/cache/stats", h.GetCacheStats)
		}
	}

	return router
}