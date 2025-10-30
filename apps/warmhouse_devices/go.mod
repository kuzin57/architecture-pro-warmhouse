module github.com/warmhouse/warmhouse_devices

go 1.24.0

require (
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/lerenn/asyncapi-codegen v0.46.3
	github.com/lib/pq v1.10.9
	github.com/stretchr/testify v1.9.0
	github.com/warmhouse/libraries/convert v0.0.0-00010101000000-000000000000
	github.com/warmhouse/libraries/rabbitmq v0.0.0-00010101000000-000000000000
	github.com/warmhouse/libraries/scheduler v0.0.0-00010101000000-000000000000
	go.uber.org/fx v1.24.0
	go.uber.org/zap v1.27.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-co-op/gocron v1.37.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/dig v1.19.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/warmhouse/libraries/rabbitmq => ../libraries/rabbitmq

replace github.com/warmhouse/libraries/convert => ../libraries/convert

replace github.com/warmhouse/libraries/scheduler => ../libraries/scheduler
