package config

import "github.com/warmhouse/libraries/rabbitmq"

type Secrets struct {
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
	} `yaml:"database"`
	RabbitMQ        rabbitmq.Config `yaml:"rabbitmq"`
	JWTSecret       string          `yaml:"jwt_secret"`
	SmarthomeAPIKey string          `yaml:"smarthome_api_key"`
}

func NewSecrets(secretsPath string) *Secrets {
	return &Secrets{}
}
