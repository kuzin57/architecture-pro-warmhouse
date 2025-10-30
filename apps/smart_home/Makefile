# Makefile for smart_home

.PHONY: generate clean build test

# Generate AsyncAPI and OpenAPI code
generate:
	@echo "Generating AsyncAPI code..."
	@mkdir -p generated/async
	@rm -rf generated/async/*.go
	asyncapi-codegen -i ../../schemas/asyncapi.yaml -o generated/async/asyncapi.gen.go -p async
	@echo "Code generation completed!"
	
	@echo "Generating OpenAPI client..."
	@mkdir -p generated/client
	@rm -rf generated/client/*.go
	oapi-codegen -config codegen-config.yaml ../../schemas/openapi.yaml
	@echo "OpenAPI client generation completed!"

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	rm -rf generated/async/*.go
	rm -rf generated/client/*.go
	@echo "Clean completed!"

# Build the application
build: generate
	@echo "Building application..."
	go build -o bin/smarthome .
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
	./bin/smarthome

# Generate and run
dev: generate
	@echo "Running in development mode..."
	go run .

# Install code generation tools
install-tools:
	@echo "Installing asyncapi-codegen..."
	go install github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@latest
	@echo "Installing oapi-codegen..."
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	@echo "Tools installed!"
