# ==== Backend ====

run:
	cd backend && go run ./cmd/server/main.go

build:
	cd backend && go build -o ../bin/app ./cmd/server/main.go

test:
	cd backend && go test -v ./...

format:
	cd backend && goimports -local github.com/aliskhannn/order-service -w .

# ==== Docker ====

docker-up:
	docker compose up --build

docker-down:
	docker compose down -v

.PHONY: producer

producer:
	docker compose exec kafka kafka-console-producer.sh --bootstrap-server kafka:9092 --topic ${TOPIC}

# ==== Frontend ====

npm-install:
	cd frontend && npm install

npm-dev:
	cd frontend && npm run dev

npm-build:
	cd frontend && npm run build
