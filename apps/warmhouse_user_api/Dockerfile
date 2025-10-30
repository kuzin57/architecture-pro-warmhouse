FROM golang:latest

WORKDIR /

COPY schemas ./schemas

WORKDIR /app

RUN go install github.com/lerenn/asyncapi-codegen/cmd/asyncapi-codegen@latest
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

COPY apps/libraries ./libraries
COPY apps/warmhouse_user_api ./warmhouse_user_api

WORKDIR /app/warmhouse_user_api

RUN cat /schemas/asyncapi.yaml

RUN go mod download
RUN go get -u github.com/lerenn/asyncapi-codegen/pkg/extensions

RUN make generate
RUN make async-api-generate

RUN go build -o warmhouse_user_api ./cmd/main

CMD ["./warmhouse_user_api"]
