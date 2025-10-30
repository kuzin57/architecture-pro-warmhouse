package config

type Config struct {
	App struct {
		Port string `yaml:"port"`
	} `yaml:"app"`
	RabbitMQ struct {
		URL string `yaml:"url"`
	} `yaml:"rabbitmq"`
}
