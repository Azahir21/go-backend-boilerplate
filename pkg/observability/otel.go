package observability

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Config holds the configuration for observability
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	TempoEndpoint  string // OTLP endpoint for Tempo (e.g., "http://localhost:4318")
	EnableTracing  bool
	EnableMetrics  bool
}

// Provider holds the OpenTelemetry providers
type Provider struct {
	TracerProvider *trace.TracerProvider
	MeterProvider  *metric.MeterProvider
	log            *logrus.Logger
}

// NewProvider initializes and returns a new OpenTelemetry provider
func NewProvider(ctx context.Context, cfg Config, log *logrus.Logger) (*Provider, error) {
	var tracerProvider *trace.TracerProvider
	var meterProvider *metric.MeterProvider
	var err error

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Setup tracing if enabled
	if cfg.EnableTracing && cfg.TempoEndpoint != "" {
		tracerProvider, err = newTracerProvider(ctx, res, cfg.TempoEndpoint)
		if err != nil {
			log.Warnf("Failed to create tracer provider: %v. Tracing will be disabled.", err)
		} else {
			otel.SetTracerProvider(tracerProvider)
			log.Info("OpenTelemetry tracing initialized successfully")
		}
	}

	// Setup metrics if enabled
	if cfg.EnableMetrics {
		meterProvider, err = newMeterProvider(ctx, res)
		if err != nil {
			log.Warnf("Failed to create meter provider: %v. Metrics collection will be disabled.", err)
		} else {
			otel.SetMeterProvider(meterProvider)
			log.Info("OpenTelemetry metrics initialized successfully")
		}
	}

	// Set global propagator for context propagation (W3C Trace Context)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Provider{
		TracerProvider: tracerProvider,
		MeterProvider:  meterProvider,
		log:            log,
	}, nil
}

// newTracerProvider creates a new trace provider with OTLP exporter
func newTracerProvider(ctx context.Context, res *resource.Resource, endpoint string) (*trace.TracerProvider, error) {
	// Create OTLP HTTP exporter for Tempo
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(), // Use insecure for local development
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP trace exporter: %w", err)
	}

	// Create trace provider with batch span processor
	tp := trace.NewTracerProvider(
		trace.WithResource(res),
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(5*time.Second),
			trace.WithMaxExportBatchSize(512),
		),
		trace.WithSampler(trace.AlwaysSample()), // Sample all traces in development
	)

	return tp, nil
}

// newMeterProvider creates a new meter provider with Prometheus exporter
func newMeterProvider(ctx context.Context, res *resource.Resource) (*metric.MeterProvider, error) {
	// Create Prometheus exporter
	exporter, err := prometheus.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus exporter: %w", err)
	}

	// Create meter provider with Prometheus exporter
	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(exporter),
	)

	return mp, nil
}

// Shutdown gracefully shuts down the OpenTelemetry providers
func (p *Provider) Shutdown(ctx context.Context) error {
	var err error

	if p.TracerProvider != nil {
		if shutdownErr := p.TracerProvider.Shutdown(ctx); shutdownErr != nil {
			err = fmt.Errorf("failed to shutdown tracer provider: %w", shutdownErr)
			p.log.Error(err)
		} else {
			p.log.Info("Tracer provider shut down successfully")
		}
	}

	if p.MeterProvider != nil {
		if shutdownErr := p.MeterProvider.Shutdown(ctx); shutdownErr != nil {
			if err != nil {
				err = fmt.Errorf("%w; failed to shutdown meter provider: %v", err, shutdownErr)
			} else {
				err = fmt.Errorf("failed to shutdown meter provider: %w", shutdownErr)
			}
			p.log.Error(err)
		} else {
			p.log.Info("Meter provider shut down successfully")
		}
	}

	return err
}
