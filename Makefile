.PHONY: build run test lint fmt tidy up down

build:
	go build ./...

run:
	go run ./cmd/server

test:
	go test ./...

lint:
	go vet ./...

fmt:
	gofmt -l -w .

tidy:
	go mod tidy

up:
	docker compose up --build

down:
	docker compose down
