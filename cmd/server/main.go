package main

import (
	"context"
	"errors"
	"github.com/aliskhannn/order-service/internal/api"
	"github.com/aliskhannn/order-service/internal/config"
	"github.com/aliskhannn/order-service/internal/repository"
	"github.com/aliskhannn/order-service/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("error creating connection pool: %v", err)
	}

	repo := repository.NewOrderRepo(dbpool)
	orderService := service.NewOrderService(repo)
	orderHTTPHandler := api.NewOrderHTTPHandler(orderService)

	r := chi.NewRouter()
	r.Get("/order/{id}", orderHTTPHandler.GetOrderByID)

	server := &http.Server{
		Addr:    cfg.Server.HTTPPort,
		Handler: r,
	}

	go func() {
		log.Printf("HTTP server listening on %s", cfg.Server.HTTPPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down HTTP server...")
	if err = server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("could not shutdown HTTP server: %v", err)
	}

	if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
		log.Fatalln("timeout exceeded, forcing shutdown")
	}

	log.Println("closing database pool...")
	dbpool.Close()
}
