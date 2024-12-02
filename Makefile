.PHONY: build run docs tidy swag install-linter lint

build:
	go build -o build/bin cmd/main.go

run: 
	build build/bin

docs:
	swag init -g ./cmd/main.go -o ./docs --parseDependency --parseInternal

tidy:
	go mod tidy

swag:
	swag init --generalInfo ./cmd/main.go --output ./docs

install-linter:
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint:
	golangci-lint run