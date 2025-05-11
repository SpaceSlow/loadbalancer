package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/SpaceSlow/loadbalancer/config"
	clientsrepo "github.com/SpaceSlow/loadbalancer/internal/repository/clients"
	"github.com/SpaceSlow/loadbalancer/internal/service/clients"
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

	repo, err := clientsrepo.NewPostgresRepo(ctx, cfg.DB)
	if err != nil {
		slog.Error("Connect to db error occurred", slog.String("error", err.Error()))
		return
	}
	service := clients.NewService(repo)

	srv, err := server.NewServer(ctx, cfg, service)
	if err != nil {
		slog.Error("Initialize server error occurred", slog.String("error", err.Error()))
		return
	}

	if err = srv.Start(); err != nil {
		slog.Error("Server failed", slog.String("error", err.Error()))
	}
}
