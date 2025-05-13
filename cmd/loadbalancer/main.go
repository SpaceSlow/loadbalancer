package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"golang.org/x/sync/errgroup"

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

	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	g, ctx := errgroup.WithContext(rootCtx)
	context.AfterFunc(ctx, func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.MaxTimeoutShutdown)
		defer cancel()

		<-timeoutCtx.Done()
		slog.Error("Failed to gracefully shutdown the service")
	})

	slog.Info("Initialize repo")
	repo, err := clientsrepo.NewPostgresRepo(ctx, cfg.DB)
	if err != nil {
		slog.Error("Initialize repo error occurred", slog.String("error", err.Error()))
		return
	}
	g.Go(func() error {
		defer slog.Info("Closed clients repo")
		<-ctx.Done()
		repo.Close()
		return nil
	})

	limiter := ratelimiter.NewRateLimiter(ctx, &cfg.RateLimiter)

	service := clients.NewService(cfg.RateLimiter.DefaultBucket, repo, limiter)

	mux := router.NewRouter(service)

	slog.Info("Initialize balancer")
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

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Balancer.Port),
		Handler: mux,
	}

	g.Go(func() error {
		slog.Info("Start balancer", slog.String("address", server.Addr))
		return server.ListenAndServe()
	})

	g.Go(func() error {
		defer slog.Info("Stopped http server")
		<-ctx.Done()
		shutdownTimeoutCtx, cancelShutdownTimeoutCtx := context.WithTimeout(
			context.Background(),
			cfg.HTTPServerTimeoutShutdown,
		)
		defer cancelShutdownTimeoutCtx()
		return server.Shutdown(shutdownTimeoutCtx)
	})
	err = g.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("Error occurred", slog.String("error", err.Error()))
		return
	}
	slog.Info("Gracefully shutdown all services")
}
