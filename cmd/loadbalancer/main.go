package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/SpaceSlow/loadbalancer/config"
	clientsrepo "github.com/SpaceSlow/loadbalancer/internal/repository/clients"
	"github.com/SpaceSlow/loadbalancer/internal/service/clients"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/balancer"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/ratelimiter"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/router"
)

func main() {
	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(slog.New(handler))

	path := flag.String("config", filepath.Join("config", "config.yaml"), "path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(*path)
	if err != nil {
		slog.Error("Load config error occurred", slog.String("error", err.Error()))
		return
	}

	ctx := context.Background()

	repo, err := clientsrepo.NewPostgresRepo(ctx, cfg.DB)
	if err != nil {
		slog.Error("Connect to db error occurred", slog.String("error", err.Error()))
		return
	}

	limiter := ratelimiter.NewRateLimiter(ctx, &cfg.RateLimiter)

	service := clients.NewService(repo, limiter)

	mux := router.NewRouter(service)

	b, err := balancer.NewBalancer(ctx, &cfg.Balancer)
	if err != nil {
		slog.Error("Create balancer error occurred", slog.String("error", err.Error()))
		return
	}

	clientsSlice, err := service.List(ctx)
	if err == nil {
		for _, c := range clientsSlice {
			limiter.AddBucket(c.ID, c.Capacity, c.RPS)
		}
	}

	mux.Handle("/", limiter.Middleware(b.Handler()))

	slog.Info("Starting balancer", slog.Int("port", cfg.Balancer.Port))
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Balancer.Port),
		Handler: mux,
	}
	if err = server.ListenAndServe(); err != nil {
		slog.Error("Server failed", slog.String("error", err.Error()))
	}
}
