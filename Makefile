.DEFAULT := build

.PHONY: fmt vet build test run coverage swagger goose

test:
	go test ./...

coverage: test
	go test -cover ./...

fmt: coverage
	go fmt ./...

vet: fmt
	golangci-lint run

swagger:
	swag init -d cmd,inner --parseInternal

goose:
	goose up -dir ./migrations

build: vet goose swagger
	go build cmd/main.go

run: goose swagger
	go run cmd/main.go