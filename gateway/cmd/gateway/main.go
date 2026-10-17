package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/cache"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/config"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/metrics"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/ratelimit"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/server"
	"github.com/Jiangge89/mini-exchange-backend/gateway/internal/upstream"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	entries := make([]upstream.Entry, len(cfg.Upstreams))
	for i, u := range cfg.Upstreams {
		entries[i] = upstream.Entry{Name: u.Name, URL: u.URL}
	}
	pool := upstream.NewPool(entries)
	selector := upstream.NewSelector(pool)
	healthChecker := upstream.NewHealthChecker(pool)
	cacheStore := cache.New(cfg.RedisURL)
	limiter := ratelimit.New(cfg.RedisURL)
	metrics.Register()

	go healthChecker.Start(context.Background())

	srv := server.New(cfg, selector, cacheStore, limiter)

	go func() {
		slog.Info("gateway starting", "addr", cfg.Addr())
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown error", "error", err)
	}
}
