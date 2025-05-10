package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/SpaceSlow/loadbalancer/config"
	"github.com/SpaceSlow/loadbalancer/internal/transport/http/server"
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
	s, err := server.NewServer(ctx, cfg, nil) // TODO add client service

	if err = s.Start(); err != nil {
		slog.Error("Server failed", slog.String("error", err.Error()))
	}
}
