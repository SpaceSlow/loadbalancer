package balancer

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"

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

	backends []*Backend
}

func NewRoundRobinBalancer(cfg *config.BalancerConfig) (*RoundRobinBalancer, error) {
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

func (b *RoundRobinBalancer) Handler(w http.ResponseWriter, r *http.Request) {
	backend := b.nextAvailableBackend()
	if backend == nil {
		slog.Error("No available backends", slog.String("uri", r.RequestURI))
		http.Error(w, "No available backends", http.StatusServiceUnavailable)
		return
	}
	backend.Proxy.ErrorHandler = backend.ProxyErrorHandler()

	backend.Proxy.ServeHTTP(w, r)
}
