package balancer

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/SpaceSlow/loadbalancer/config"
	"github.com/SpaceSlow/loadbalancer/pkg/networks"
	"github.com/SpaceSlow/loadbalancer/pkg/statuscode"
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

type RoundRobinBalancer struct {
	counter     atomic.Uint64
	backendsNum uint64

	port     uint16
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
		port:        uint16(cfg.Port),
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

func (b *RoundRobinBalancer) Start() error {
	slog.Info("Starting balancer", slog.Int("port", int(b.port)))
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", b.port),
		Handler: http.HandlerFunc(b.Handler),
	}
	return server.ListenAndServe()
}

func (b *RoundRobinBalancer) Handler(w http.ResponseWriter, r *http.Request) {
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
}
