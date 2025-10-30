package config

type Config struct {
	UserAPI struct {
		URL string `yaml:"url"`
	} `yaml:"user_api"`
}
