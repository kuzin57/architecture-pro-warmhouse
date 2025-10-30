package main

import (
	"flag"

	"github.com/warmhouse/warmhouse_user_api/cmd/app"

	"go.uber.org/fx"
)

func main() {
	var confPath, secretsPath string

	flag.StringVar(&confPath, "config", "config/config.yaml", "path to config")
	flag.StringVar(&secretsPath, "secrets", "config/secrets.yaml", "path to secrets")

	flag.Parse()

	fx.New(app.Create(confPath, secretsPath)).Run()
}
