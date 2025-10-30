# Makefile for warmhouse_user_api

.PHONY: generate clean build test

# Generate OpenAPI code
generate:
	@echo "Generating OpenAPI code..."
	@mkdir -p internal/generated/server
	@rm -rf internal/generated/server/*.go
	oapi-codegen -config codegen-config.yaml api/openapi.yaml
	@echo "Code generation completed!"

async-api-generate:
	@echo "Generating AsyncAPI code..."
	@mkdir -p internal/generated/async
	@rm -rf internal/generated/async/*.go
	asyncapi-codegen -i ../../schemas/asyncapi.yaml -o internal/generated/async/asyncapi.gen.go -p async
	@echo "Code generation completed!"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf internal/generated/server/*.go
	rm -rf internal/generated/async/*.go
	@echo "Clean completed!"

# Build the application
build: generate
	@echo "Building application..."
	go build -o bin/warmhouse_user_api ./cmd/main
	@echo "Build completed!"

# Run tests
test:
	@echo "Running tests..."
	go test ./...
	@echo "Tests completed!"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download
	@echo "Dependencies installed!"

# Run the application
run: build
	@echo "Running application..."
	./bin/warmhouse_user_api

# Generate and run
dev: generate
	@echo "Running in development mode..."
	go run ./cmd/main
