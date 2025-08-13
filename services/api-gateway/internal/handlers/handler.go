package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Fadil369/NPHIES/services/api-gateway/internal/auth"
	"github.com/Fadil369/NPHIES/services/api-gateway/internal/config"
	"github.com/Fadil369/NPHIES/services/api-gateway/internal/kafka"
	"github.com/Fadil369/NPHIES/services/api-gateway/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	config   *config.Config
	logger   *logrus.Logger
	db       *sql.DB
	redis    *redis.Client
	kafka    *kafka.Producer
	auth     *auth.Service
	metrics  *MetricsCollector
}

type MetricsCollector struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	ActiveRequests  prometheus.Gauge
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config, logger *logrus.Logger) (*Handler, error) {
	// Initialize database connection
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure database connection pool
	db.SetMaxOpenConns(cfg.Database.MaxConnections)
	db.SetMaxIdleConns(cfg.Database.MaxConnections / 2)
	db.SetConnMaxLifetime(time.Duration(cfg.Database.MaxIdleTime) * time.Minute)

	// Initialize Redis client
	redisOpts, err := redis.ParseURL(cfg.Redis.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}
	redisClient := redis.NewClient(redisOpts)

	if err := redisClient.Ping(redisClient.Context()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	// Initialize Kafka producer
	kafkaProducer, err := kafka.NewProducer(cfg.Kafka.Brokers, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	// Initialize auth service
	authService := auth.NewService(cfg.JWT.Secret, time.Duration(cfg.JWT.Expiration)*time.Second)

	// Initialize metrics
	metrics := &MetricsCollector{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_request_duration_seconds",
				Help: "HTTP request duration in seconds",
			},
			[]string{"method", "endpoint"},
		),
		ActiveRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_active_requests",
				Help: "Number of active HTTP requests",
			},
		),
	}

	prometheus.MustRegister(metrics.RequestsTotal)
	prometheus.MustRegister(metrics.RequestDuration)
	prometheus.MustRegister(metrics.ActiveRequests)

	return &Handler{
		config:  cfg,
		logger:  logger,
		db:      db,
		redis:   redisClient,
		kafka:   kafkaProducer,
		auth:    authService,
		metrics: metrics,
	}, nil
}

// Close closes all connections
func (h *Handler) Close() error {
	if h.db != nil {
		h.db.Close()
	}
	if h.redis != nil {
		h.redis.Close()
	}
	if h.kafka != nil {
		h.kafka.Close()
	}
	return nil
}

// Health check endpoints
// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	status := "healthy"
	statusCode := http.StatusOK

	// Check database connection
	if err := h.db.Ping(); err != nil {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
		h.logger.Errorf("Database health check failed: %v", err)
	}

	// Check Redis connection
	if err := h.redis.Ping(c.Request.Context()).Err(); err != nil {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
		h.logger.Errorf("Redis health check failed: %v", err)
	}

	c.JSON(statusCode, gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"service":   "api-gateway",
		"version":   "1.0.0",
	})
}

// ReadinessCheck godoc
// @Summary Readiness check
// @Description Check if the service is ready to serve requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /ready [get]
func (h *Handler) ReadinessCheck(c *gin.Context) {
	ready := true
	dependencies := make(map[string]string)

	// Check database
	if err := h.db.Ping(); err != nil {
		ready = false
		dependencies["database"] = "not ready"
	} else {
		dependencies["database"] = "ready"
	}

	// Check Redis
	if err := h.redis.Ping(c.Request.Context()).Err(); err != nil {
		ready = false
		dependencies["redis"] = "not ready"
	} else {
		dependencies["redis"] = "ready"
	}

	// Check Kafka (simplified check)
	dependencies["kafka"] = "ready" // In production, implement proper Kafka health check

	statusCode := http.StatusOK
	if !ready {
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"ready":        ready,
		"dependencies": dependencies,
		"timestamp":    time.Now().UTC(),
	})
}

// MetricsHandler serves Prometheus metrics
func (h *Handler) MetricsHandler(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

// Authentication endpoints
// GetToken godoc
// @Summary Get authentication token
// @Description Authenticate user and return JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginRequest true "User credentials"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/token [post]
func (h *Handler) GetToken(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request format",
			Message: err.Error(),
		})
		return
	}

	// TODO: Implement proper authentication with OAuth2/OIDC
	// For now, using a simple mock implementation
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Missing credentials",
			Message: "Username and password are required",
		})
		return
	}

	// Generate JWT token
	token, err := h.auth.GenerateToken(req.Username, []string{"read", "write"})
	if err != nil {
		h.logger.Errorf("Failed to generate token: %v", err)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Token generation failed",
			Message: "Unable to generate authentication token",
		})
		return
	}

	// Log successful authentication
	h.logAuditEvent("auth.login", req.Username, c.ClientIP(), map[string]interface{}{
		"username": req.Username,
		"success":  true,
	})

	c.JSON(http.StatusOK, models.TokenResponse{
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    h.config.JWT.Expiration,
		RefreshToken: "", // TODO: Implement refresh tokens
		Scope:        "read write",
	})
}

// RefreshToken godoc
// @Summary Refresh authentication token
// @Description Refresh an existing JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param token body models.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} models.TokenResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	// TODO: Implement refresh token logic
	c.JSON(http.StatusNotImplemented, models.ErrorResponse{
		Error:   "Not implemented",
		Message: "Refresh token functionality is not yet implemented",
	})
}

// logAuditEvent logs an audit event to Kafka
func (h *Handler) logAuditEvent(eventType, userID, clientIP string, data map[string]interface{}) {
	auditEvent := map[string]interface{}{
		"eventId":     uuid.New().String(),
		"eventType":   eventType,
		"userId":      userID,
		"clientIP":    clientIP,
		"timestamp":   time.Now().UTC(),
		"service":     "api-gateway",
		"data":        data,
	}

	eventData, err := json.Marshal(auditEvent)
	if err != nil {
		h.logger.Errorf("Failed to marshal audit event: %v", err)
		return
	}

	if err := h.kafka.Publish(h.config.Kafka.Topics.AuditTrail, string(eventData)); err != nil {
		h.logger.Errorf("Failed to publish audit event: %v", err)
	}
}