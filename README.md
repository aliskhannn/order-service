# Order Service

## Overview

This project is an order service with a backend written in Go and a frontend using Vite.  
The backend uses Kafka for messaging and `go-cache` package for caching.

## Project Structure

```
.
├── cmd/              # Entry points (Kafka consumer, HTTP server)
├── config/           # Configuration files
├── internal/         # Application internal packages
│   ├── api/          # HTTP handlers, routers, server
│   ├── config/       # Config parsing logic
│   ├── infra/        # Infrastructure (cache, kafka)
│   ├── kafka/        # Kafka message handlers
│   ├── logger/       # Logger setup
│   ├── model/        # Data models
│   ├── repository/   # Database repositories
│   ├── service/      # Business logic
│   └── validator/    # Input validation
├── migrations/       # DB migrations
├── website/          # Frontend application
├── .env.example
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── Makefile
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