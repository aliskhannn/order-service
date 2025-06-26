package main

import (
	"context"
	"github.com/aliskhannn/order-service/internal/config"
	"github.com/aliskhannn/order-service/internal/infra/kafka"
	"github.com/aliskhannn/order-service/internal/kafka/handlers"
	"github.com/aliskhannn/order-service/internal/repository"
	"github.com/aliskhannn/order-service/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os/signal"
	"sync"
	"syscall"
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
	orderService := service.NewOrderService(repo)

	orderCreatedTopicHandler := handlers.NewOrderCreatedTopicHandler(orderService)

	consumer := kafka.NewConsumer(cfg.Kafka.Topic, cfg.Kafka.Addr, orderCreatedTopicHandler)
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
