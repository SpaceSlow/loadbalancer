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
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/balancer"
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

	b, err := balancer.NewBalancer(context.Background(), &cfg.Balancer)
	if err != nil {
		slog.Error("Create balancer error occurred", slog.String("error", err.Error()))
		return
	}

	slog.Info("Starting balancer", slog.Int("port", cfg.Balancer.Port))
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Balancer.Port),
		Handler: b.Handler(),
	}

	if err = server.ListenAndServe(); err != nil {
		slog.Error("Server failed", slog.String("error", err.Error()))
	}
}
