package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Fadil369/NPHIES/services/terminology-service/internal/config"
	"github.com/Fadil369/NPHIES/services/terminology-service/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Setup logger
	setupLogger(cfg.LogLevel)

	logrus.Info("Starting NPHIES Terminology Service...")

	// Setup Gin router
	r := gin.Default()

	// Initialize handlers
	handler := handlers.NewHandler(cfg)

	// Setup routes
	setupRoutes(r, handler)

	// Start server
	logrus.Infof("Server starting on port %s", cfg.Port)
	
	// Setup graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(":" + cfg.Port); err != nil {
			logrus.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for shutdown signal
	<-quit
	logrus.Info("Shutting down NPHIES Terminology Service...")
}

func setupLogger(level string) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	
	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func setupRoutes(r *gin.Engine, handler *handlers.Handler) {
	// Health check endpoints
	r.GET("/health", handler.Health)
	r.GET("/ready", handler.Ready)

	// API routes
	api := r.Group("/api/v1")
	{
		// Code systems management
		codeSystems := api.Group("/code-systems")
		{
			codeSystems.GET("", handler.ListCodeSystems)
			codeSystems.POST("", handler.CreateCodeSystem)
			codeSystems.GET("/:id", handler.GetCodeSystem)
			codeSystems.PUT("/:id", handler.UpdateCodeSystem)
			codeSystems.DELETE("/:id", handler.DeleteCodeSystem)
		}

		// Code lookup and mapping
		codes := api.Group("/codes")
		{
			codes.GET("/lookup/:system/:code", handler.LookupCode)
			codes.POST("/validate", handler.ValidateCode)
			codes.POST("/map", handler.MapCode)
		}

		// Value sets management
		valueSets := api.Group("/value-sets")
		{
			valueSets.GET("", handler.ListValueSets)
			valueSets.POST("", handler.CreateValueSet)
			valueSets.GET("/:id", handler.GetValueSet)
			valueSets.PUT("/:id", handler.UpdateValueSet)
			valueSets.DELETE("/:id", handler.DeleteValueSet)
		}

		// Terminology administration
		admin := api.Group("/admin")
		{
			admin.POST("/import", handler.ImportTerminology)
			admin.POST("/cache/refresh", handler.RefreshCache)
			admin.GET("/stats", handler.GetStatistics)
		}
	}
}