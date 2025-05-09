package balancer

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

	"github.com/SpaceSlow/loadbalancer/config"
	"github.com/SpaceSlow/loadbalancer/pkg/networks"
	"github.com/SpaceSlow/loadbalancer/pkg/statuscode"
)

type RoundRobinBalancer struct {
	counter     atomic.Uint64
	backendsNum uint64

	backends []*Backend
}

func NewRoundRobinBalancer(ctx context.Context, cfg *config.BalancerConfig) (*RoundRobinBalancer, error) {
	backends := make([]*Backend, 0, len(cfg.Backends))

	for _, backendCfg := range cfg.Backends {
		backendURL, err := url.Parse(backendCfg.URL)
		if err != nil {
			return nil, fmt.Errorf("[config] incorrect url: %w", err)
		}

		backend := &Backend{
			Proxy: httputil.NewSingleHostReverseProxy(backendURL),
			URL:   backendURL,
		}
		backend.SetAlive(true)
		backends = append(backends, backend)

		go backend.HealthCheckLoop(ctx, &backendCfg.HealthCheck)
	}

	return &RoundRobinBalancer{
		backendsNum: uint64(len(backends)),
		backends:    backends,
	}, nil
}

func (b *RoundRobinBalancer) nextAvailableBackend() *Backend {
	for range b.backendsNum {
		next := b.counter.Add(1)
		backend := b.backends[next%b.backendsNum]
		if backend.IsAlive() {
			return backend
		}
		slog.Warn("Skip unavailable backend", slog.String("backend", backend.URL.String()))
	}
	return nil
}

func (b *RoundRobinBalancer) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		backend := b.nextAvailableBackend()
		if backend == nil {
			slog.Error("No available backends", slog.String("uri", r.RequestURI))
			http.Error(w, "No available backends", http.StatusServiceUnavailable)
			return
		}
		backend.Proxy.ErrorHandler = backend.ProxyErrorHandler()

		lw := statuscode.NewResponseWriter(w)
		backend.Proxy.ServeHTTP(lw, r)
		slog.Info(
			"Request",
			slog.String("ip", networks.ParseIP(r.RemoteAddr)),
			slog.String("method", r.Method),
			slog.String("uri", r.RequestURI),
			slog.String("user-agent", r.UserAgent()),
			slog.Int("status_code", lw.StatusCode),
			slog.String("backend", backend.URL.String()),
		)
	})
}
