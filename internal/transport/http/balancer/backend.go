package balancer

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/SpaceSlow/loadbalancer/config"
)

type Backend struct {
	Proxy *httputil.ReverseProxy
	URL   *url.URL
	alive atomic.Bool
}

func (b *Backend) IsAlive() bool {
	return b.alive.Load()
}

func (b *Backend) SetAlive(isAlive bool) {
	b.alive.Store(isAlive)
}

func (b *Backend) HealthCheckLoop(ctx context.Context, cfg *config.HealthCheckConfig) {
	ticker := time.NewTicker(cfg.Interval)
	defer ticker.Stop()

	uri, err := url.JoinPath(b.URL.String(), cfg.Path)
	if err != nil {
		slog.Error("Healthcheck url with path join error", slog.String("error", err.Error()))
		return
	}

	for {
		select {
		case <-ctx.Done():
			slog.Info("Healthcheck stopped", slog.String("backend", b.URL.String()))
		case <-ticker.C:
			response, err := http.Get(uri)
			if err != nil {
				b.SetAlive(false)
				slog.Error("Healthcheck error occurred", slog.String("error", err.Error()))
				continue
			}
			b.SetAlive(response.StatusCode == http.StatusOK)
		}
	}
}

func (b *Backend) ProxyErrorHandler() func(w http.ResponseWriter, req *http.Request, err error) {
	return func(w http.ResponseWriter, req *http.Request, err error) {
		slog.Error(
			"Backend unavailable",
			slog.String("backend", b.URL.String()),
			slog.String("uri", req.RequestURI),
			slog.String("error", err.Error()),
		)
		b.SetAlive(false)
		http.Error(w, "Backend unavailable", http.StatusServiceUnavailable)
	}
}
