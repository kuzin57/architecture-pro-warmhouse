package config

type Config struct {
	App struct {
		Port string `yaml:"port"`
	} `yaml:"app"`
}
