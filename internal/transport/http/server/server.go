package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SpaceSlow/loadbalancer/config"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/balancer"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/clients"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/ratelimiter"
)

type Server struct {
	server         *http.Server
	clientHandlers clients.Handlers
}

func NewServer(ctx context.Context, cfg *config.Config, clientService clients.Service) (*Server, error) {
	const basePathLayout = "/api/v1%s"

	mux := http.NewServeMux()

	clientHandlers := clients.NewHandlers(clientService)

	mux.Handle(fmt.Sprintf(basePathLayout, "/clients/"), clientHandlers.Clients())
	mux.Handle(fmt.Sprintf(basePathLayout, "/clients/{id}/"), clientHandlers.ClientByID())

	b, err := balancer.NewBalancer(ctx, &cfg.Balancer)
	if err != nil {
		slog.Error("Create balancer error occurred", slog.String("error", err.Error()))
		return nil, err
	}
	mux.Handle("/", b.Handler())

	limiter := ratelimiter.NewRateLimiter(ctx, &cfg.RateLimiter)

	slog.Info("Starting balancer", slog.Int("port", cfg.Balancer.Port))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Balancer.Port),
		Handler: limiter.Middleware(mux),
	}

	return &Server{server: server}, nil
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}
