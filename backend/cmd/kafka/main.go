package main

import (
	"context"
	"go.uber.org/zap"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aliskhannn/order-service/internal/config"
	"github.com/aliskhannn/order-service/internal/infra/kafka"
	"github.com/aliskhannn/order-service/internal/kafka/handlers/order"
	"github.com/aliskhannn/order-service/internal/logger"
	orderrepo "github.com/aliskhannn/order-service/internal/repository/order"
	ordersvc "github.com/aliskhannn/order-service/internal/service/order"
	"github.com/aliskhannn/order-service/internal/validator"
)

func main() {
	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()
	log := logger.CreateLogger(cfg.Env)

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal("error creating connection pool", zap.Error(err))
	}

	repo := orderrepo.New(log, dbpool)
	orderService := ordersvc.New(log, nil, repo)
	val := validator.New()

	orderCreatedHandler := order.NewCreateHandler(log, val, orderService)

	consumer := kafka.NewConsumer(cfg.Kafka.GroupID, cfg.Kafka.Topic, cfg.Kafka.Brokers, log, orderCreatedHandler)
	wg.Add(1)
	go consumer.ConsumeMessage(ctx, &wg)

	log.Info("Kafka consumer started...")

	<-ctx.Done()
	log.Info("shutdown signal received")

	wg.Wait()

	log.Info("closing kafka consumer...")
	err = consumer.Close()
	if err != nil {
		log.Error("error closing kafka consumer: %v", zap.Error(err))
	}

	log.Info("closing database pool...")
	dbpool.Close()
}
