# Order Service

## Overview

This project is an order service with a backend written in Go and a frontend using Vite. The backend uses Kafka for messaging and go-cache package for caching.

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/aliskhannn/order-service.git
cd order-service
````

### 2. Create `.env` file

Copy the example environment file and fill in your configuration values:

```bash
cp .env.example .env
```

Edit `.env` as needed.

---

## Running the Application with Docker Compose

To build and start all services (backend, frontend, Kafka etc):

```bash
make docker-up
```

To stop and remove all containers, networks, and volumes created by Docker Compose:

```bash
make docker-down
```

---

## Running Tests

Run all backend tests with:

```bash
make test
```

---

## Frontend

After running `make docker-up`, open your browser at:

```
http://localhost:8081/
```

to access the frontend application.

---

## Useful Makefile Commands

### Backend

* `make run` — run the Go server
* `make build` — build the Go server binary
* `make test` — run tests
* `make format` — format Go code with `goimports`

### Docker

* `make docker-up` — build and start all Docker containers
* `make docker-down` — stop and clean up Docker containers and volumes

### Frontend

* `make npm-install` — install frontend dependencies
* `make npm-dev` — run frontend in development mode
* `make npm-build` — build frontend for production

---

## Notes

* Make sure Docker and Docker Compose are installed on your machine.
* The `.env` file is required for proper configuration.
* Backend runs on port `8080`.
* Frontend is served on port `8081`.
