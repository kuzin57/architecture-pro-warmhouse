package config

type Secrets struct {
	JWTSecret string `yaml:"jwt_secret"`
}

func NewSecrets(secretsPath string) *Secrets {
	return &Secrets{}
}
