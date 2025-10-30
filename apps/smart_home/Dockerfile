# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies for code generation tools
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Install code generation tools
RUN go install github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@latest
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# Copy libraries first (needed for go.mod replace directives)
COPY apps/libraries ./libraries

# Copy smart_home source code
COPY apps/smart_home ./smart_home

# Copy schemas for code generation
COPY schemas ./schemas

# Download dependencies
RUN cd smart_home && go mod download

# Generate code
RUN cd smart_home && mkdir -p generated/async generated/client
RUN cd smart_home && asyncapi-codegen -i ../../schemas/asyncapi.yaml -o generated/async/asyncapi.gen.go -p async || true
RUN cd smart_home && oapi-codegen -config codegen-config.yaml ../../schemas/openapi.yaml || true

# Build the application
RUN cd smart_home && CGO_ENABLED=0 GOOS=linux go build -o smarthome .

# Runtime stage
FROM alpine:latest

# Install necessary runtime dependencies (including postgresql-client for migrations)
RUN apk --no-cache add ca-certificates tzdata postgresql-client

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/smart_home/smarthome .

# Copy generated code (if generation worked)
COPY --from=builder /app/smart_home/generated ./generated

# Copy config files
COPY --from=builder /app/smart_home/secrets.yaml .
COPY --from=builder /app/smart_home/config.yaml .

# Copy entrypoint script from builder stage
COPY --from=builder /app/smart_home/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose the port the app runs on
EXPOSE 8080

# Use entrypoint to run migrations on startup
ENTRYPOINT ["/entrypoint.sh"]
CMD ["./smarthome"]