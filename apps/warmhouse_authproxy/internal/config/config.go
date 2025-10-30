package config

type Config struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
	JwtDurationMinutes int `yaml:"jwt_duration_minutes"`
}

func NewConfig(configPath string) *Config {
	return &Config{}
}
