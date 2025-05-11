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
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/dto"
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
			b.healthCheck(uri)
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
		dto.WriteErrorResponse(w, http.StatusServiceUnavailable, "Backend unavailable")
	}
}

func (b *Backend) healthCheck(uri string) {
	response, err := http.Get(uri)

	var errAttr slog.Attr
	switch {
	case err != nil && b.IsAlive():
		errAttr = slog.String("error", err.Error())
		fallthrough
	case response != nil && response.StatusCode != http.StatusOK && b.IsAlive():
		b.SetAlive(false)
		slog.Error(
			"Backend has become unavailable (healthcheck)",
			slog.String("backend", b.URL.String()),
			errAttr,
		)
	case err == nil && response.StatusCode == http.StatusOK && !b.IsAlive():
		b.SetAlive(true)
		slog.Info(
			"Backend has become available (healthcheck)",
			slog.String("backend", b.URL.String()),
		)
	}
}
