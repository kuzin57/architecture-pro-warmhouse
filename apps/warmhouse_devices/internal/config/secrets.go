package config

import "github.com/warmhouse/libraries/rabbitmq"

type Secrets struct {
	Pg struct {
		DSN string `yaml:"dsn"`
	} `yaml:"pg"`
	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`
}
