# Order Service

## Overview

This project is an order service with a backend written in Go and a frontend using Vite.  
The backend uses Kafka for messaging and `go-cache` package for caching.

## Project Structure

```
.
├── backend
│   ├── cmd
│   │   ├── kafka
│   │   │   └── main.go           # Kafka consumer service entry point
│   │   └── server
│   │       └── main.go           # HTTP server entry point
│   ├── config
│   │   └── config.yml            # Configuration file
│   ├── internal
│   │   ├── api
│   │   │   ├── handlers
│   │   │   │   └── order
│   │   │   │       ├── get_handler.go        # Handler to get order by ID
│   │   │   │       └── get_handler_test.go  # Unit tests for handler
│   │   │   │   └── response.go               # Helper functions to write JSON/error responses
│   │   │   ├── router
│   │   │   │   └── router.go      # HTTP router setup
│   │   │   └── server
│   │   │       └── server.go      # HTTP server initialization
│   │   ├── config
│   │   │   └── config.go          # Configuration parsing logic
│   │   ├── infra
│   │   │   ├── cache
│   │   │   │   ├── cache.go       # Caching logic
│   │   │   │   └── errors.go      # Custom cache errors
│   │   │   └── kafka
│   │   │       ├── consumer.go    # Kafka consumer client
│   │   │       └── kafka.go       # Kafka producer/consumer setup
│   │   ├── kafka
│   │   │   └── handlers/order
│   │   │       ├── order_created_handler.go       # Handle "order created" messages
│   │   │       └── errors.go                      # Custom errors for Kafka handlers
│   │   ├── logger
│   │   │   └── logger.go          # Zap logger setup
│   │   ├── model
│   │   │   ├── delivery.go
│   │   │   ├── item.go
│   │   │   ├── order.go
│   │   │   └── payment.go
│   │   ├── repository
│   │   │   └── order
│   │   │       ├── order_repo.go
│   │   │       └── errors.go
│   │   ├── service
│   │   │   └── order
│   │   │       ├── order_service.go
│   │   │       └── order_service_test.go
│   │   └── validator
│   │       └── validator.go       # Input validation logic
│   ├── migrations                 # Database migrations
│   ├── Dockerfile                 # Backend Dockerfile
│   ├── go.mod
│   └── go.sum
├── frontend
│   ├── src
│   │   ├── app/
│   │   ├── entities/
│   │   ├── features/
│   │   └── pages/
│   └── Dockerfile
├── .env.example                   # Example environment variables
├── docker-compose.yml             # Docker Compose configuration
├── Makefile                       # Helper commands
└── README.md
```

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/aliskhannn/order-service.git
cd order-service
````

### 2. Create `.env` file

```bash
cp .env.example .env
```

Edit `.env` as needed.

---

## Running the Application with Docker Compose

```bash
make docker-up      # Build and start all services
make docker-down    # Stop and clean up containers and volumes
```

## Kafka Producer

Start a Kafka console producer to send messages to a topic:

```bash
make producer TOPIC=order.created
```

Then you can type messages in the console and they will be sent to the Kafka topic.

---

## Running Tests

```bash
make test           # Run all backend tests
```

---

## Useful Makefile Commands

### Backend

* `make test` — run backend tests
* `make format` — format Go code
* `make lint` — run linters

### Docker

* `make docker-up` — build and start all Docker containers
* `make docker-down` — stop and clean up Docker containers

### Frontend

* `make npm-install` — install frontend dependencies
* `make npm-dev` — run frontend in development mode
* `make npm-build` — build frontend for production