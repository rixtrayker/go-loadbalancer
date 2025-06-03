# Docker Monitoring Setup

This directory contains Docker configurations for comprehensive monitoring of the Go Load Balancer using Prometheus and Grafana.

## Directory Structure

```
docker/
├── prometheus/
│   ├── prometheus.yml      # Prometheus configuration
│   ├── alert_rules.yml     # Alert rules for load balancer
│   └── alertmanager.yml    # AlertManager configuration
├── grafana/
│   ├── datasources.yml     # Grafana datasources (Prometheus + InfluxDB)
│   ├── dashboards.yml      # Dashboard provisioning
│   └── load-balancer-dashboard.json  # Pre-built dashboard
└── README.md               # This file
```

## Services Overview

### Prometheus Stack
- **Prometheus** (`:9090`) - Metrics collection and storage
- **AlertManager** (`:9093`) - Alert handling and notifications
- **Node Exporter** (`:9100`) - System metrics collection

### Grafana Instances
- **Standard Grafana** (`:3000`) - InfluxDB datasource for k6 metrics
- **Enhanced Grafana** (`:3001`) - Prometheus + InfluxDB datasources

## Quick Start

### 1. Start Complete Monitoring Stack
```bash
make k6-monitor-full
```

### 2. Start Prometheus-only Stack
```bash
make k6-monitor-prometheus
```

### 3. Start InfluxDB-only Stack (Original)
```bash
make k6-monitor
```

## Available Endpoints

| Service | URL | Description |
|---------|-----|-------------|
| Load Balancer | http://localhost:8080 | Your Go load balancer |
| Grafana (InfluxDB) | http://localhost:3000 | K6 metrics visualization |
| Grafana (Prometheus) | http://localhost:3001 | System + app metrics |
| Prometheus | http://localhost:9090 | Metrics database |
| InfluxDB | http://localhost:8086 | K6 time-series data |
| AlertManager | http://localhost:9093 | Alert management |
| Node Exporter | http://localhost:9100 | System metrics |

## Running Tests

### With InfluxDB Output (Original)
```bash
make k6-stress K6_VUS=50 K6_DURATION=5m
```

### With Prometheus Output
```bash
make k6-stress-prometheus K6_VUS=50 K6_DURATION=5m
```

### With Both Outputs
```bash
# Start full stack
make k6-monitor-full

# Run test (outputs to both InfluxDB and Prometheus)
make k6-stress-prometheus
```

## Monitoring Features

### Prometheus Metrics
- Load balancer health and performance
- Backend service availability
- System metrics (CPU, memory, disk)
- K6 test metrics
- Custom application metrics

### Grafana Dashboards
- **Load Balancer Overview** - Service health, request rates, response times
- **K6 Test Results** - Test execution metrics and performance
- **System Metrics** - Server resource utilization
- **Alerts Dashboard** - Active alerts and their status

### Alert Rules
- Load balancer downtime
- Backend service failures
- High response times
- Error rate spikes
- K6 test failures

## Configuration

### Customizing Prometheus
Edit `prometheus/prometheus.yml` to:
- Add new scrape targets
- Modify scrape intervals
- Configure recording rules

### Customizing Alerts
Edit `prometheus/alert_rules.yml` to:
- Add new alert conditions
- Modify thresholds
- Create custom alert groups

### Customizing Grafana
- Datasources: `grafana/datasources.yml`
- Dashboard provisioning: `grafana/dashboards.yml`
- Custom dashboards: Add JSON files to `grafana/` directory

## Docker Compose Files

### Main Stack (`docker-compose.yml`)
- Load balancer and backends
- K6 testing
- InfluxDB + Grafana

### Monitoring Stack (`docker-compose.monitoring.yml`)
- Prometheus
- AlertManager
- Node Exporter
- Enhanced Grafana

## Troubleshooting

### Services Not Starting
```bash
# Check logs
docker-compose logs prometheus
docker-compose logs grafana-enhanced

# Validate configurations
docker-compose config
```

### Prometheus Not Scraping
1. Check `prometheus/prometheus.yml` syntax
2. Verify target endpoints are accessible
3. Check Prometheus targets page: http://localhost:9090/targets

### Grafana Connection Issues
1. Verify datasource configuration in `grafana/datasources.yml`
2. Test connection in Grafana UI
3. Check network connectivity between containers

## Best Practices

1. **Resource Limits**: Set appropriate CPU/memory limits in production
2. **Data Retention**: Configure Prometheus retention policies
3. **Security**: Use authentication and HTTPS in production
4. **Backup**: Regular backup of Grafana dashboards and Prometheus data
5. **Monitoring**: Monitor the monitoring stack itself

## Production Considerations

- Use external storage for Prometheus data
- Configure proper authentication for Grafana
- Set up SSL/TLS termination
- Implement proper backup strategies
- Configure log rotation
- Set resource limits and health checks 