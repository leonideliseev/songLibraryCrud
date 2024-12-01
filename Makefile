.PHONY: build
build:
	go build -o build/bin cmd/main.go

.PHONY: run
run: build
	build/bin

.PHONY: docs
docs:
	swag init -g ./cmd/main.go -o ./docs --parseDependency --parseInternal