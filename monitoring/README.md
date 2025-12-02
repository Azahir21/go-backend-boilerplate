# Monitoring Stack

This directory contains the complete observability stack for the Go Backend Boilerplate.

## Quick Start

```bash
# Start the monitoring stack
docker compose up -d

# Check all services are running
docker compose ps

# View logs
docker compose logs -f

# Stop the monitoring stack
docker compose down

# Clean up (including volumes)
docker compose down -v
```

## Services

| Service | Port | URL | Credentials |
|---------|------|-----|-------------|
| Grafana | 3000 | http://localhost:3000 | admin / admin |
| Prometheus | 9090 | http://localhost:9090 | - |
| Alertmanager | 9093 | http://localhost:9093 | - |
| Loki | 3100 | http://localhost:3100 | - |
| Tempo | 3200, 4317, 4318 | http://localhost:3200 | - |

## Architecture

```
Go App → Prometheus (metrics)
       → Tempo (traces via OTLP)
       → Promtail → Loki (logs)

All visualized in → Grafana
```

## Prerequisites

1. **Docker Network**: The monitoring stack expects the `go-backend-network` to exist. This is created by the main application's docker-compose.

2. **Application Configuration**: The Go backend must have the following environment variable set:
   ```bash
   TEMPO_ENDPOINT=localhost:4318
   ```

## Starting Everything

From the repository root:

```bash
# 1. Start the main application (creates go-backend-network)
docker compose up -d

# 2. Start the monitoring stack
cd monitoring
docker compose up -d

# 3. Verify everything is running
docker compose ps
cd ..
docker compose ps
```

## Accessing Dashboards

1. **Grafana**: Open http://localhost:3000
   - Login: admin / admin
   - Go to Dashboards → Browse → "Go Backend - Overview"

2. **Prometheus**: Open http://localhost:9090
   - Try query: `rate(http_requests_total[5m])`

3. **Check Metrics**: Open http://localhost:8080/metrics
   - Should show Prometheus metrics from your application

## Troubleshooting

### Services won't start

```bash
# Check logs
docker compose logs <service-name>

# Common issues:
# 1. Port already in use - check with: netstat -an | grep <port>
# 2. Network doesn't exist - start main application first
```

### No metrics in Grafana

1. Check Prometheus targets: http://localhost:9090/targets
2. Ensure "go-backend" target is UP
3. Check application is running: `curl http://localhost:8080/metrics`

### No logs in Loki

1. Check Promtail logs: `docker compose logs promtail`
2. Verify Loki is running: `curl http://localhost:3100/ready`
3. Check Docker socket permissions

### No traces in Tempo

1. Verify application environment: `echo $TEMPO_ENDPOINT`
2. Check Tempo logs: `docker compose logs tempo`
3. Test OTLP endpoint: `curl http://localhost:4318/v1/traces`

## Configuration Files

- `docker-compose.yaml` - All monitoring services
- `prometheus.yml` - Prometheus configuration and scrape targets
- `loki.yaml` - Loki configuration
- `promtail.yaml` - Log collection configuration
- `tempo.yaml` - Trace storage configuration
- `alertmanager.yml` - Alert routing configuration
- `alert.rules.yml` - Alert definitions
- `grafana/provisioning/` - Grafana data sources and dashboards

## Customization

### Add New Alert

Edit `alert.rules.yml`:

```yaml
- alert: MyNewAlert
  expr: my_metric > threshold
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Alert description"
```

Then restart Prometheus:
```bash
docker compose restart prometheus
```

### Add New Dashboard

1. Create dashboard in Grafana UI
2. Export as JSON
3. Save to `grafana/dashboards/my-dashboard.json`
4. Restart Grafana: `docker compose restart grafana`

### Configure Alert Notifications

Edit `alertmanager.yml` to add Slack, email, or other receivers.

## Resource Usage

Approximate resource usage for all monitoring services:
- **Memory**: ~1-2 GB
- **CPU**: Low (increases with scale)
- **Disk**: Depends on retention settings

To reduce resource usage:
- Lower Prometheus scrape interval
- Reduce retention periods
- Limit log ingestion rate

## More Information

See the comprehensive guide: [../docs/monitoring.md](../docs/monitoring.md)
