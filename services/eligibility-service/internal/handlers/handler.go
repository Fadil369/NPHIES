package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Fadil369/NPHIES/services/eligibility-service/internal/cache"
	"github.com/Fadil369/NPHIES/services/eligibility-service/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	config    *config.Config
	logger    *logrus.Logger
	db        *sql.DB
	redis     *redis.Client
	cache     *cache.Manager
	kafka     *kafka.Writer
	metrics   *MetricsCollector
	startTime time.Time
}

type MetricsCollector struct {
	RequestsTotal        *prometheus.CounterVec
	RequestDuration      *prometheus.HistogramVec
	CacheHits            prometheus.Counter
	CacheMisses          prometheus.Counter
	EligibilityChecks    prometheus.Counter
	DatabaseQueries      *prometheus.CounterVec
	ActiveConnections    prometheus.Gauge
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

	// Initialize cache manager
	cacheManager := cache.NewManager(redisClient, time.Duration(cfg.Business.CacheTTL)*time.Second)

	// Initialize Kafka writer
	kafkaWriter := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Brokers...),
		Topic:    cfg.Kafka.Topics.AuditTrail,
		Balancer: &kafka.LeastBytes{},
	}

	// Initialize metrics
	metrics := &MetricsCollector{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "eligibility_requests_total",
				Help: "Total number of eligibility requests",
			},
			[]string{"method", "endpoint", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "eligibility_request_duration_seconds",
				Help: "Eligibility request duration in seconds",
			},
			[]string{"method", "endpoint"},
		),
		CacheHits: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "eligibility_cache_hits_total",
				Help: "Total number of cache hits",
			},
		),
		CacheMisses: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "eligibility_cache_misses_total",
				Help: "Total number of cache misses",
			},
		),
		EligibilityChecks: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "eligibility_checks_total",
				Help: "Total number of eligibility checks performed",
			},
		),
		DatabaseQueries: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "database_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"operation"},
		),
		ActiveConnections: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "database_active_connections",
				Help: "Number of active database connections",
			},
		),
	}

	prometheus.MustRegister(metrics.RequestsTotal)
	prometheus.MustRegister(metrics.RequestDuration)
	prometheus.MustRegister(metrics.CacheHits)
	prometheus.MustRegister(metrics.CacheMisses)
	prometheus.MustRegister(metrics.EligibilityChecks)
	prometheus.MustRegister(metrics.DatabaseQueries)
	prometheus.MustRegister(metrics.ActiveConnections)

	return &Handler{
		config:    cfg,
		logger:    logger,
		db:        db,
		redis:     redisClient,
		cache:     cacheManager,
		kafka:     kafkaWriter,
		metrics:   metrics,
		startTime: time.Now(),
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

// HealthCheck godoc
// @Summary Health check
// @Description Check if the eligibility service is healthy
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	status := "healthy"
	statusCode := http.StatusOK
	checks := make(map[string]string)

	// Check database connection
	if err := h.db.Ping(); err != nil {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
		checks["database"] = "unhealthy: " + err.Error()
		h.logger.Errorf("Database health check failed: %v", err)
	} else {
		checks["database"] = "healthy"
	}

	// Check Redis connection
	if err := h.redis.Ping(c.Request.Context()).Err(); err != nil {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
		checks["redis"] = "unhealthy: " + err.Error()
		h.logger.Errorf("Redis health check failed: %v", err)
	} else {
		checks["redis"] = "healthy"
	}

	c.JSON(statusCode, gin.H{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"service":   "eligibility-service",
		"version":   "1.0.0",
		"checks":    checks,
		"uptime":    time.Since(h.startTime).String(),
	})
}

// ReadinessCheck godoc
// @Summary Readiness check
// @Description Check if the eligibility service is ready to serve requests
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
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

	// Check if we can query essential data
	var count int
	err := h.db.QueryRow("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'").Scan(&count)
	if err != nil {
		ready = false
		dependencies["schema"] = "not ready"
	} else {
		dependencies["schema"] = "ready"
	}

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

// logAuditEvent logs an audit event to Kafka
func (h *Handler) logAuditEvent(ctx context.Context, eventType, userID, clientIP string, data map[string]interface{}) {
	auditEvent := map[string]interface{}{
		"eventId":     uuid.New().String(),
		"eventType":   eventType,
		"userId":      userID,
		"clientIP":    clientIP,
		"timestamp":   time.Now().UTC(),
		"service":     "eligibility-service",
		"data":        data,
	}

	eventData, err := json.Marshal(auditEvent)
	if err != nil {
		h.logger.Errorf("Failed to marshal audit event: %v", err)
		return
	}

	message := kafka.Message{
		Key:   []byte(eventType),
		Value: eventData,
		Time:  time.Now(),
	}

	if err := h.kafka.WriteMessages(ctx, message); err != nil {
		h.logger.Errorf("Failed to publish audit event: %v", err)
	}
}