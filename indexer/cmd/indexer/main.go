package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jiangge89/mini-exchange-backend/indexer/internal/api"
	"github.com/Jiangge89/mini-exchange-backend/indexer/internal/chain"
	"github.com/Jiangge89/mini-exchange-backend/indexer/internal/config"
	"github.com/Jiangge89/mini-exchange-backend/indexer/internal/scanner"
	"github.com/Jiangge89/mini-exchange-backend/indexer/internal/store"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := store.NewPostgres(context.Background(), cfg.DatabaseURL())
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	ethClient, err := chain.NewClient(cfg.GatewayURL)
	if err != nil {
		slog.Error("failed to connect to gateway", "error", err)
		os.Exit(1)
	}

	sc := scanner.New(ethClient, db, cfg.TokenContractAddress)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sc.Run(ctx)

	srv := api.NewServer(cfg.Addr(), db)
	go func() {
		slog.Info("indexer API starting", "addr", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("API server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down")

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	srv.Shutdown(shutdownCtx)
}
