package main

import (
	"context"
	"errors"
	"github.com/aliskhannn/order-service/internal/api/handlers/order"
	"github.com/aliskhannn/order-service/internal/api/router"
	server2 "github.com/aliskhannn/order-service/internal/api/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"go.uber.org/zap"

	"github.com/aliskhannn/order-service/internal/config"
	gocache "github.com/aliskhannn/order-service/internal/infra/cache"
	"github.com/aliskhannn/order-service/internal/logger"
	orderrepo "github.com/aliskhannn/order-service/internal/repository/order"
	ordersvc "github.com/aliskhannn/order-service/internal/service/order"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := logger.CreateLogger(cfg.Env)

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal("error creating connection pool", zap.Error(err))
	}

	repo := orderrepo.New(log, dbpool)
	cache := gocache.New(cfg.Cache.DefaultExpiration, cfg.Cache.CleanupInterval, log, repo)

	// Preload cache with existing orders
	err = cache.Preload(ctx, cfg.Cache.PreloadLimit)
	if err != nil {
		log.Fatal("error preloading cache", zap.Error(err))
	}

	orderService := ordersvc.New(log, cache, repo)
	orderGetHandler := order.NewGetHandler(log, orderService)

	r := router.New(orderGetHandler)
	server := server2.New(cfg.Server.HTTPPort, r)

	go func() {
		log.Info("starting HTTP server", zap.String("port", cfg.Server.HTTPPort))
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("shutting down HTTP server...")
	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Error("could not shutdown HTTP server", zap.Error(err))
		os.Exit(1)
	}

	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		log.Fatal("timeout exceeded, forcing shutdown")
	}

	log.Info("closing database pool...")
	dbpool.Close()
}
