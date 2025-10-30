module github.com/warmhouse/warmhouse_user_api

go 1.24.0

require (
	github.com/go-chi/chi/v5 v5.2.3
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lerenn/asyncapi-codegen v0.46.3
	github.com/lib/pq v1.10.9
	github.com/oapi-codegen/runtime v1.1.2
	github.com/stretchr/testify v1.11.1
	github.com/warmhouse/libraries/convert v0.0.0-00010101000000-000000000000
	github.com/warmhouse/libraries/rabbitmq v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.24.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.43.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/warmhouse/libraries/convert => ../libraries/convert

replace github.com/warmhouse/libraries/rabbitmq => ../libraries/rabbitmq
