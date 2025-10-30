FROM golang:latest

WORKDIR /app

COPY libraries ./libraries

COPY warmhouse_devices ./warmhouse_devices

RUN cd warmhouse_devices && go mod download

RUN cd warmhouse_devices && go build -o warmhouse_devices ./cmd/main

WORKDIR /app/warmhouse_devices
CMD ["./warmhouse_devices"]
