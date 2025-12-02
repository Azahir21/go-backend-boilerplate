package observability

import (
	"context"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// Metrics holds custom application metrics
type Metrics struct {
	// HTTP metrics
	HTTPRequestsTotal      metric.Int64Counter
	HTTPRequestDuration    metric.Float64Histogram
	HTTPRequestsInFlight   metric.Int64UpDownCounter

	// Database metrics
	DBQueryDuration        metric.Float64Histogram
	DBConnectionsActive    metric.Int64UpDownCounter
	DBConnectionsIdle      metric.Int64UpDownCounter
	DBConnectionsMax       metric.Int64Gauge

	// Application metrics
	ServiceUptime          metric.Float64Counter
	GoroutineCount         metric.Int64Gauge

	log                    *logrus.Logger
	startTime              time.Time
}

// NewMetrics creates and initializes custom metrics
func NewMetrics(log *logrus.Logger) (*Metrics, error) {
	meter := otel.Meter("github.com/azahir21/go-backend-boilerplate")

	// HTTP metrics
	httpRequestsTotal, err := meter.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	httpRequestDuration, err := meter.Float64Histogram(
		"http_request_duration_seconds",
		metric.WithDescription("HTTP request latency in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	httpRequestsInFlight, err := meter.Int64UpDownCounter(
		"http_requests_in_flight",
		metric.WithDescription("Current number of HTTP requests being processed"),
		metric.WithUnit("{request}"),
	)
	if err != nil {
		return nil, err
	}

	// Database metrics
	dbQueryDuration, err := meter.Float64Histogram(
		"db_query_duration_seconds",
		metric.WithDescription("Database query latency in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	dbConnectionsActive, err := meter.Int64UpDownCounter(
		"db_pool_connections_active",
		metric.WithDescription("Number of active database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return nil, err
	}

	dbConnectionsIdle, err := meter.Int64UpDownCounter(
		"db_pool_connections_idle",
		metric.WithDescription("Number of idle database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return nil, err
	}

	dbConnectionsMax, err := meter.Int64Gauge(
		"db_pool_connections_max",
		metric.WithDescription("Maximum number of database connections"),
		metric.WithUnit("{connection}"),
	)
	if err != nil {
		return nil, err
	}

	// Application metrics
	serviceUptime, err := meter.Float64Counter(
		"service_uptime_seconds",
		metric.WithDescription("Service uptime in seconds"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, err
	}

	goroutineCount, err := meter.Int64Gauge(
		"goroutine_count",
		metric.WithDescription("Number of goroutines"),
		metric.WithUnit("{goroutine}"),
	)
	if err != nil {
		return nil, err
	}

	m := &Metrics{
		HTTPRequestsTotal:    httpRequestsTotal,
		HTTPRequestDuration:  httpRequestDuration,
		HTTPRequestsInFlight: httpRequestsInFlight,
		DBQueryDuration:      dbQueryDuration,
		DBConnectionsActive:  dbConnectionsActive,
		DBConnectionsIdle:    dbConnectionsIdle,
		DBConnectionsMax:     dbConnectionsMax,
		ServiceUptime:        serviceUptime,
		GoroutineCount:       goroutineCount,
		log:                  log,
		startTime:            time.Now(),
	}

	// Start background metrics collection
	go m.collectRuntimeMetrics()

	return m, nil
}

// collectRuntimeMetrics collects runtime metrics periodically
func (m *Metrics) collectRuntimeMetrics() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()

		// Collect goroutine count
		goroutines := int64(runtime.NumGoroutine())
		m.GoroutineCount.Record(ctx, goroutines)

		// Collect service uptime
		uptime := time.Since(m.startTime).Seconds()
		m.ServiceUptime.Add(ctx, uptime)
	}
}

// RecordHTTPRequest records HTTP request metrics
func (m *Metrics) RecordHTTPRequest(ctx context.Context, method, path, status string, duration float64) {
	// Record request count
	m.HTTPRequestsTotal.Add(ctx, 1,
		metric.WithAttributes(
			attribute.String("method", method),
			attribute.String("path", path),
			attribute.String("status", status),
		),
	)

	// Record request duration
	m.HTTPRequestDuration.Record(ctx, duration,
		metric.WithAttributes(
			attribute.String("method", method),
			attribute.String("path", path),
		),
	)
}

// RecordDBQuery records database query metrics
func (m *Metrics) RecordDBQuery(ctx context.Context, operation string, duration float64) {
	m.DBQueryDuration.Record(ctx, duration,
		metric.WithAttributes(
			attribute.String("operation", operation),
		),
	)
}

// UpdateDBConnectionStats updates database connection pool statistics
func (m *Metrics) UpdateDBConnectionStats(ctx context.Context, active, idle, max int64) {
	m.DBConnectionsActive.Add(ctx, active)
	m.DBConnectionsIdle.Add(ctx, idle)
	m.DBConnectionsMax.Record(ctx, max)
}
