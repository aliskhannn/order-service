run:
	go run ./cmd/server/main.go

build:
	go build -o bin/app ./cmd/server/main.go

test:
	go test -v ./...

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v

.PHONY: producer

producer:
	docker compose exec kafka kafka-console-producer.sh --bootstrap-server kafka:9092 --topic ${TOPIC}

format:
	goimports -local github.com/aliskhannn/order-service -w .
