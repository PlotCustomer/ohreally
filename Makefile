.PHONY: build test vet run tidy up down

build:
	go build ./...

test:
	go test ./...

vet:
	go vet ./...

run:
	go run ./cmd/server

tidy:
	go mod tidy

up:
	docker compose up -d postgres

down:
	docker compose down
