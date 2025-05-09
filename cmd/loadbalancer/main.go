package main

import (
	"flag"
	"log/slog"
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

	b, err := balancer.NewBalancer(&cfg.Balancer)
	if err != nil {
		slog.Error("Create balancer error occurred", slog.String("error", err.Error()))
		return
	}

	if err = b.Start(); err != nil {
		slog.Error("Server failed", slog.String("error", err.Error()))
	}
}
