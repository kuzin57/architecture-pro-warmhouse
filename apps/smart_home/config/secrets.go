package config

import "github.com/warmhouse/libraries/rabbitmq"

type Secrets struct {
	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`
	APIKey   string          `yaml:"api_key"`
}
