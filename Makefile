# ==== Backend ====

# Run all backend tests with verbose output
test:
	cd backend && go test -v ./...

# Format Go code using goimports
format:
	cd backend && goimports -local github.com/aliskhannn/order-service -w .

# Run linters: vet + golangci-lint
lint:
	cd backend && go vet ./... && golangci-lint run ./...

# ==== Docker ====

# Build and start all Docker services
docker-up:
	docker compose up --build

# Stop and remove all Docker services and volumes
docker-down:
	docker compose down -v

.PHONY: producer

# Start a Kafka console producer to send messages to a topic
producer:
	docker compose exec kafka kafka-console-producer.sh --bootstrap-server kafka:9092 --topic ${TOPIC}

# ==== Frontend ====

# Install frontend dependencies
npm-install:
	cd frontend && npm install

# Start frontend in development mode
npm-dev:
	cd frontend && npm run dev

# Build frontend for production
npm-build:
	cd frontend && npm run build