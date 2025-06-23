run:
	go run ./cmd/app/main.go

build:
	go build -o bin/app ./cmd/app/main.go

test:
	go test -v ./...

docker-up:
	docker compose up

docker-down:
	docker compose down