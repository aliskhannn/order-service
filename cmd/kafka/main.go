package main

import (
	"context"
	"github.com/aliskhannn/order-service/internal/config"
	cache2 "github.com/aliskhannn/order-service/internal/infra/cache"
	"github.com/aliskhannn/order-service/internal/infra/kafka"
	kafkahandlers "github.com/aliskhannn/order-service/internal/kafka/handlers"
	"github.com/aliskhannn/order-service/internal/repository"
	"github.com/aliskhannn/order-service/internal/service"
	"github.com/aliskhannn/order-service/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	var wg sync.WaitGroup

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.MustLoad()

	dbpool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatalf("error creating connection pool: %v", err)
	}

	repo := repository.NewOrderRepo(dbpool)
	cache := cache2.NewGoCache(5*time.Minute, 10*time.Minute)
	orderService := service.NewOrderService(repo, cache)
	val := validator.New()

	orderCreatedHandler := kafkahandlers.NewOrderCreatedHandler(orderService, val)

	consumer := kafka.NewConsumer(cfg.Kafka.GroupID, cfg.Kafka.Topic, cfg.Kafka.Addr, orderCreatedHandler)
	wg.Add(1)
	go consumer.ConsumeMessage(ctx, &wg)

	log.Println("Kafka consumer started...")

	<-ctx.Done()
	log.Println("shutdown signal received")

	wg.Wait()

	log.Println("closing kafka consumer...")
	err = consumer.Close()
	if err != nil {
		log.Printf("error closing kafka consumer: %v", err)
	}

	log.Println("closing database pool...")
	dbpool.Close()
}
