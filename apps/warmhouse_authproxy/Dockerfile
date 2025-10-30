FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY warmhouse_authproxy/go.mod warmhouse_authproxy/go.sum ./
RUN go mod download

COPY warmhouse_authproxy/cmd ./cmd
COPY warmhouse_authproxy/internal ./internal
COPY warmhouse_authproxy/config ./config

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o warmhouse_authproxy ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/warmhouse_authproxy .

COPY --from=builder /app/config ./config

EXPOSE 8087

CMD ["./warmhouse_authproxy", "-next-host", "warmhouse_user_api:8080"]
