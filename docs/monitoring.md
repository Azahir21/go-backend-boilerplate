# Monitoring & Observability Guide

This guide covers the full observability stack implemented for the Go Backend Boilerplate, including metrics, logs, tracing, and alerting.

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Components](#components)
- [Accessing Dashboards](#accessing-dashboards)
- [Available Metrics](#available-metrics)
- [Adding Custom Metrics](#adding-custom-metrics)
- [Tracing](#tracing)
- [Logs](#logs)
- [Alerts](#alerts)
- [Troubleshooting](#troubleshooting)

## Overview

The monitoring stack provides comprehensive observability for the Go backend application using:

- **Prometheus** - Metrics collection and storage
- **Loki** - Log aggregation
- **Tempo** - Distributed tracing
- **Grafana** - Unified dashboard and visualization
- **Alertmanager** - Alert routing and notifications
- **OpenTelemetry** - Instrumentation SDK for metrics and traces

## Architecture

```
┌──────────────────┐      scrape       ┌───────────────────┐
│    Go Backend     │ ───────────────► │    Prometheus      │
│  (OTEL + metrics) │                  │ (metrics storage)  │
└──────────────────┘                   └─────────┬─────────┘
        │ logs via stdout                             │
        ▼                                             ▼
┌───────────────────┐                         ┌──────────────────┐
│     Promtail      │                         │      Loki        │
│ (log forwarding)  │ ───────────────────────►│ (log database)   │
└───────────────────┘                         └────────┬─────────┘
                                                       │
       ┌───────────────────────────────────────────────┘
       ▼
┌───────────────────────────┐
│           Grafana         │
│ Dashboards + Alerts + UI  │
└───────────────────────────┘

Tracing: Go Backend → OTLP HTTP → Tempo → Grafana
```

## Quick Start

### Prerequisites

- Docker and Docker Compose installed
- Go backend application running (or ready to start)
- Sufficient disk space for metrics/logs/traces storage

### Starting the Monitoring Stack

1. **Start the monitoring services:**

```bash
cd monitoring
docker compose up -d
```

2. **Verify all services are running:**

```bash
docker compose ps
```

You should see all services in a "running" state:
- prometheus
- loki
- promtail
- tempo
- grafana
- alertmanager

3. **Start your Go backend application:**

```bash
# From the root of the repository
docker compose up -d app

# Or run locally
go run cmd/main.go
```

4. **Verify metrics are being collected:**

Visit `http://localhost:8080/metrics` - you should see Prometheus metrics.

### Accessing Dashboards

Once everything is running, access the following UIs:

| Service | URL | Credentials |
|---------|-----|-------------|
| **Grafana** | http://localhost:3000 | admin / admin |
| **Prometheus** | http://localhost:9090 | None |
| **Alertmanager** | http://localhost:9093 | None |

## Components

### Prometheus (Port 9090)

Prometheus scrapes metrics from the Go application every 15 seconds via the `/metrics` endpoint.

**Configuration:** `monitoring/prometheus.yml`

**Key Features:**
- Automatic service discovery
- Time-series data storage
- PromQL query language
- Alert rule evaluation

**Common Queries:**

```promql
# Request rate
rate(http_requests_total[5m])

# P95 latency
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))

# Error rate
sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))
```

### Loki (Port 3100)

Loki aggregates logs from Docker containers via Promtail.

**Configuration:** `monitoring/loki.yaml`

**Features:**
- 31-day log retention
- Efficient log storage
- Label-based log querying
- Integration with Grafana

**Log Query Examples:**

```logql
# All logs from go-backend
{container="go-backend-app"}

# Error logs only
{container="go-backend-app"} |= "error"

# JSON field extraction
{container="go-backend-app"} | json | level="error"
```

### Tempo (Port 3200, 4317, 4318)

Tempo stores distributed traces from the application.

**Configuration:** `monitoring/tempo.yaml`

**Features:**
- OTLP gRPC (4317) and HTTP (4318) receivers
- Trace storage and querying
- Integration with Grafana and Loki

**Accessing Traces:**
1. Go to Grafana → Explore
2. Select "Tempo" data source
3. Search by Trace ID or use Service Graph

### Grafana (Port 3000)

Unified visualization platform for metrics, logs, and traces.

**Configuration:** 
- Data sources: `monitoring/grafana/provisioning/datasources/datasources.yaml`
- Dashboards: `monitoring/grafana/dashboards/`

**Pre-configured Dashboards:**
1. **Go Backend - Overview**: Main dashboard showing HTTP metrics, latency, errors, goroutines, memory, and database performance

**Navigation:**
- Home → Dashboards → Browse → "Go Backend - Overview"
- Explore → Select data source (Prometheus, Loki, or Tempo)

### Alertmanager (Port 9093)

Routes and manages alerts from Prometheus.

**Configuration:** `monitoring/alertmanager.yml`

**Alert Channels:**
To enable notifications, edit `monitoring/alertmanager.yml` and configure:
- Slack webhooks
- Email SMTP
- PagerDuty
- Webhook endpoints

**Example Slack Configuration:**

```yaml
receivers:
  - name: 'critical'
    slack_configs:
      - api_url: 'YOUR_SLACK_WEBHOOK_URL'
        channel: '#critical-alerts'
        title: 'Critical Alert: {{ .CommonLabels.alertname }}'
        text: '{{ range .Alerts }}{{ .Annotations.description }}{{ end }}'
```

### Promtail

Collects logs from Docker containers and forwards them to Loki.

**Configuration:** `monitoring/promtail.yaml`

**Features:**
- Automatic container log discovery
- JSON log parsing
- Label extraction
- Real-time log forwarding

## Available Metrics

### HTTP Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `http_requests_total` | Counter | Total HTTP requests by method, path, status |
| `http_request_duration_seconds` | Histogram | HTTP request latency distribution |
| `http_requests_in_flight` | Gauge | Current number of requests being processed |

### Database Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `db_query_duration_seconds` | Histogram | Database query latency |
| `db_pool_connections_active` | Gauge | Active database connections |
| `db_pool_connections_idle` | Gauge | Idle database connections |
| `db_pool_connections_max` | Gauge | Maximum database connections |

### Application Metrics

| Metric | Type | Description |
|--------|------|-------------|
| `service_uptime_seconds` | Counter | Service uptime in seconds |
| `goroutine_count` | Gauge | Number of goroutines |
| `go_memstats_alloc_bytes` | Gauge | Allocated memory (Go runtime) |
| `go_memstats_sys_bytes` | Gauge | System memory (Go runtime) |

## Adding Custom Metrics

### Step 1: Add Metric to the Metrics Struct

Edit `pkg/observability/metrics.go`:

```go
// Add your new metric
MyCustomCounter metric.Int64Counter
```

### Step 2: Initialize the Metric

In the `NewMetrics` function:

```go
myCustomCounter, err := meter.Int64Counter(
    "my_custom_counter",
    metric.WithDescription("Description of my counter"),
    metric.WithUnit("{unit}"),
)
if err != nil {
    return nil, err
}
```

### Step 3: Add a Recording Method

```go
func (m *Metrics) RecordMyCustomMetric(ctx context.Context, value int64) {
    m.MyCustomCounter.Add(ctx, value,
        metric.WithAttributes(
            metric.String("label_name", "label_value"),
        ),
    )
}
```

### Step 4: Use in Your Code

```go
import "github.com/azahir21/go-backend-boilerplate/pkg/observability"

// In your handler or service
func MyHandler(c *gin.Context) {
    ctx := c.Request.Context()
    
    // Record your custom metric
    if metrics := getMetricsFromContext(ctx); metrics != nil {
        metrics.RecordMyCustomMetric(ctx, 1)
    }
}
```

### Step 5: Query in Prometheus/Grafana

After deploying, your metric will be available as `my_custom_counter` in Prometheus.

## Tracing

The application uses OpenTelemetry for distributed tracing.

### Automatic Instrumentation

HTTP requests are automatically traced via the `otelgin` middleware:
- Span created for each HTTP request
- Request method, path, status code captured
- Parent-child span relationships maintained

### Viewing Traces

1. Make a request to your application
2. Go to Grafana → Explore
3. Select "Tempo" data source
4. Search by Trace ID (found in logs) or use the Service Graph
5. Click on a trace to see the full span tree

### Adding Custom Spans

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
)

func MyFunction(ctx context.Context) error {
    tracer := otel.Tracer("my-service")
    ctx, span := tracer.Start(ctx, "my-operation")
    defer span.End()
    
    // Add attributes
    span.SetAttributes(
        attribute.String("user.id", "12345"),
        attribute.Int("items.count", 42),
    )
    
    // Your code here
    
    return nil
}
```

### Connecting Traces to Logs

Logs automatically include trace IDs when using structured logging with the trace context.

## Logs

### Log Format

The application outputs structured JSON logs:

```json
{
  "time": "2024-12-02T10:30:00Z",
  "level": "info",
  "msg": "Request processed",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "latency": 0.042
}
```

### Viewing Logs in Grafana

1. Go to Grafana → Explore
2. Select "Loki" data source
3. Use LogQL to query logs:

```logql
# All application logs
{container="go-backend-app"}

# Error logs in the last hour
{container="go-backend-app"} |= "error" | json | level="error"

# Logs for a specific endpoint
{container="go-backend-app"} | json | path="/api/v1/users"

# Logs with trace correlation
{container="go-backend-app"} | json | trace_id="abc123"
```

### Log Levels

- `debug` - Detailed debugging information
- `info` - General informational messages
- `warn` - Warning messages
- `error` - Error messages
- `fatal` - Critical errors that cause shutdown

## Alerts

Pre-configured alerts in `monitoring/alert.rules.yml`:

### Active Alerts

| Alert | Severity | Condition | Duration |
|-------|----------|-----------|----------|
| **HighErrorRate** | Warning | Error rate > 5% | 5 minutes |
| **HighLatency** | Warning | P95 latency > 1s | 5 minutes |
| **ServiceDown** | Critical | Service unreachable | 1 minute |
| **HighMemoryUsage** | Warning | Memory usage > 90% | 5 minutes |
| **TooManyGoroutines** | Warning | Goroutines > 1000 | 5 minutes |
| **DatabaseConnectionPoolExhausted** | Warning | No idle connections | 2 minutes |
| **HighDatabaseLatency** | Warning | P95 DB latency > 0.5s | 5 minutes |

### Testing Alerts

To test if alerts are working:

1. **Trigger an alert condition** (e.g., stop the application for ServiceDown)
2. Check Prometheus → Alerts (http://localhost:9090/alerts)
3. Check Alertmanager → Alerts (http://localhost:9093/#/alerts)
4. If configured, check your notification channel (Slack, email, etc.)

### Adding Custom Alerts

Edit `monitoring/alert.rules.yml`:

```yaml
- alert: MyCustomAlert
  expr: my_custom_metric > 100
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Custom alert triggered"
    description: "My metric is {{ $value }}"
```

Reload Prometheus configuration:

```bash
docker compose -f monitoring/docker-compose.yaml restart prometheus
```

## Troubleshooting

### Metrics Not Appearing

1. **Check /metrics endpoint:**
   ```bash
   curl http://localhost:8080/metrics
   ```

2. **Check Prometheus targets:**
   - Go to http://localhost:9090/targets
   - Ensure "go-backend" target is UP

3. **Check Prometheus logs:**
   ```bash
   docker compose -f monitoring/docker-compose.yaml logs prometheus
   ```

### Traces Not Showing

1. **Verify Tempo is running:**
   ```bash
   docker compose -f monitoring/docker-compose.yaml ps tempo
   ```

2. **Check environment variable:**
   ```bash
   echo $TEMPO_ENDPOINT
   # Should be: localhost:4318
   ```

3. **Check application logs for OTLP errors**

### No Logs in Loki

1. **Check Promtail is running:**
   ```bash
   docker compose -f monitoring/docker-compose.yaml logs promtail
   ```

2. **Verify Docker socket access:**
   ```bash
   docker compose -f monitoring/docker-compose.yaml exec promtail ls -l /var/run/docker.sock
   ```

3. **Check Loki ingestion:**
   ```bash
   curl http://localhost:3100/ready
   ```

### Grafana Dashboards Not Loading

1. **Check Grafana data sources:**
   - Go to Grafana → Configuration → Data Sources
   - Test each data source connection

2. **Check provisioning:**
   ```bash
   docker compose -f monitoring/docker-compose.yaml logs grafana | grep -i "provisioning"
   ```

### High Resource Usage

If monitoring services consume too much resources:

1. **Adjust Prometheus retention:**
   Edit `monitoring/docker-compose.yaml`:
   ```yaml
   command:
     - '--storage.tsdb.retention.time=7d'  # Reduce from default
   ```

2. **Reduce scrape frequency:**
   Edit `monitoring/prometheus.yml`:
   ```yaml
   global:
     scrape_interval: 30s  # Increase from 15s
   ```

3. **Limit Loki retention:**
   Edit `monitoring/loki.yaml`:
   ```yaml
   limits_config:
     retention_period: 168h  # 7 days instead of 31
   ```

## Environment Variables

Configure observability via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `TEMPO_ENDPOINT` | `localhost:4318` | Tempo OTLP HTTP endpoint |
| `OTEL_SERVICE_NAME` | `go-backend-boilerplate` | Service name for traces |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | - | Override OTLP endpoint |

## Production Considerations

When deploying to production:

1. **Enable TLS** for Tempo OTLP endpoint
2. **Configure persistent volumes** for data retention
3. **Set up proper retention policies** based on your needs
4. **Configure alert notifications** (Slack, PagerDuty, etc.)
5. **Implement sampling** for high-traffic services
6. **Use service meshes** (Istio, Linkerd) for additional telemetry
7. **Consider managed services** (Grafana Cloud, DataDog, New Relic)
8. **Implement log rotation** to prevent disk exhaustion
9. **Set resource limits** for monitoring containers
10. **Regularly review and update alert thresholds**

## Next Steps

1. Customize dashboards for your specific use cases
2. Add application-specific metrics
3. Configure alert notifications
4. Implement log aggregation from multiple services
5. Set up long-term storage for metrics/logs/traces
6. Create runbooks for common alerts
7. Implement synthetic monitoring and health checks

## Resources

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [OpenTelemetry Go Documentation](https://opentelemetry.io/docs/instrumentation/go/)
- [Loki Documentation](https://grafana.com/docs/loki/)
- [Tempo Documentation](https://grafana.com/docs/tempo/)

---

**Questions or Issues?** Open an issue on GitHub or contact the DevOps team.
