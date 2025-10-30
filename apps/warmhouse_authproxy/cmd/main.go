package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/warmhouse/warmhouse_authproxy/internal/config"
	"github.com/warmhouse/warmhouse_authproxy/internal/server"
	"github.com/warmhouse/warmhouse_authproxy/internal/services/auth"
	"gopkg.in/yaml.v2"
)

func mustLoadConfig(configPath string) *config.Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	conf := &config.Config{}
	if err := yaml.Unmarshal(data, conf); err != nil {
		panic(err)
	}
	return conf
}

func mustLoadSecrets(secretsPath string) *config.Secrets {
	data, err := os.ReadFile(secretsPath)
	if err != nil {
		log.Fatalf("failed to load secrets: %v", err)
	}
	secrets := &config.Secrets{}
	if err := yaml.Unmarshal(data, secrets); err != nil {
		panic(err)
	}
	return secrets
}

func main() {
	var confPath, secretsPath, nextHost string

	flag.StringVar(&confPath, "config", "config/config.yaml", "path to config")
	flag.StringVar(&secretsPath, "secrets", "config/secrets.yaml", "path to secrets")
	flag.StringVar(&nextHost, "next-host", "localhost:8080", "path to next host")
	flag.Parse()

	var (
		conf        = mustLoadConfig(confPath)
		secrets     = mustLoadSecrets(secretsPath)
		authService = auth.NewService(secrets, conf)
		server      = server.NewServer(nextHost, conf, authService)
		ctx, cancel = context.WithCancel(context.Background())
	)

	defer cancel()

	go func() {
		if err := server.Start(); err != nil {
			log.Println("failed to start server", err)
		}
	}()

	<-ctx.Done()

	if err := server.Stop(ctx); err != nil {
		log.Println("failed to stop server", err)
	}
}
