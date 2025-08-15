.PHONY: build

build:
	@go build -o bin/nox ./cmd/nox/main.go

.PHONY: run
run: build
	@./bin/nox