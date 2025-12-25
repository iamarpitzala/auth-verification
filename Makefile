.PHONY: build run test clean dev

# Build the application
build:
	go build -o bin/auth-backend cmd/server/main.go

# Run the application
run: build
	./bin/auth-backend

# Run in development mode with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy
	go mod download

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run