package healthcheck

import (
	"context"
	"time"

	"github.com/rixtrayker/go-loadbalancer/configs"
	"github.com/rixtrayker/go-loadbalancer/internal/healthcheck/probes"
	"github.com/rixtrayker/go-loadbalancer/internal/serverpool"
	"github.com/rixtrayker/go-loadbalancer/pkg/logging"
)

// HealthChecker monitors backend health
type HealthChecker struct {
	pools      map[string]*serverpool.Pool
	configs    map[string]configs.HealthCheckConfig
	probes     map[string]probes.Probe
	logger     *logging.Logger
	stopCh     chan struct{}
	cancelFunc context.CancelFunc
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(
	pools map[string]*serverpool.Pool,
	configs map[string]configs.HealthCheckConfig,
	logger *logging.Logger,
) *HealthChecker {
	return &HealthChecker{
		pools:   pools,
		configs: configs,
		probes:  make(map[string]probes.Probe),
		logger:  logger,
		stopCh:  make(chan struct{}),
	}
}

// Start begins health checking
func (hc *HealthChecker) Start(ctx context.Context) {
	ctx, hc.cancelFunc = context.WithCancel(ctx)

	// Create probes for each backend
	for poolName, pool := range hc.pools {
		config, ok := hc.configs[poolName]
		if !ok {
			hc.logger.Warn("No health check config for pool", "pool", poolName)
			continue
		}

		for _, backend := range pool.Backends {
			var probe probes.Probe
			switch {
			case config.Path != "":
				probe = probes.NewHTTPProbe(backend.URL, config.Path, config.Method, config.Timeout)
			default:
				probe = probes.NewTCPProbe(backend.URL, config.Timeout)
			}

			backendURL := backend.URL.String()
			hc.probes[backendURL] = probe

			// Start health check for this backend
			go hc.checkBackend(ctx, pool, backend.URL.String(), config.Interval)
		}
	}
}

// Stop stops health checking
func (hc *HealthChecker) Stop() {
	if hc.cancelFunc != nil {
		hc.cancelFunc()
	}
	close(hc.stopCh)
}

// checkBackend periodically checks a backend's health
func (hc *HealthChecker) checkBackend(ctx context.Context, pool *serverpool.Pool, backendURL string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	probe, ok := hc.probes[backendURL]
	if !ok {
		hc.logger.Error("No probe found for backend", "backend", backendURL)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			healthy := probe.Check()
			pool.MarkBackendStatus(backendURL, healthy)

			if !healthy {
				hc.logger.Warn("Backend is unhealthy", "backend", backendURL)
			}
		}
	}
}
