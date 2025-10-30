# Makefile for warmhouse_devices

.PHONY: generate clean build test gates-client-generate async-api-generate

# Generate AsyncAPI code
async-api-generate:
	@echo "Generating AsyncAPI code..."
	@mkdir -p internal/generated/async
	@rm -rf internal/generated/async/*.go
	asyncapi-codegen -i ../../schemas/asyncapi.yaml -o internal/generated/async/asyncapi.gen.go -p async
	@echo "AsyncAPI code generation completed!"

# Generate all code
generate: async-api-generate
	@echo "All code generation completed!"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf internal/generated/async/*.go
	@echo "Clean completed!"

# Build the application
build: generate
	@echo "Building application..."
	go build -o bin/warmhouse_devices ./cmd/main
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
	./bin/warmhouse_devices

# Generate and run
dev: generate
	@echo "Running in development mode..."
	go run ./cmd/main
