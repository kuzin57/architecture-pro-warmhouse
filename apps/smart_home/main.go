package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	userapi "smarthome/clients/user_api"
	"smarthome/config"
	"smarthome/db"
	"smarthome/generated/async"
	"smarthome/handlers"
	"smarthome/services"

	"github.com/warmhouse/libraries/rabbitmq"
	"gopkg.in/yaml.v2"

	"github.com/gin-gonic/gin"
)

func mustLoadSecrets(secretsPath string) *config.Secrets {
	data, err := os.ReadFile(secretsPath)
	if err != nil {
		log.Fatalf("Unable to load secrets: %v\n", err)
	}

	secrets := &config.Secrets{}
	if err := yaml.Unmarshal(data, secrets); err != nil {
		log.Fatalf("Unable to unmarshal secrets: %v\n", err)
	}

	return secrets
}

func mustLoadConfig(configPath string) *config.Config {
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Unable to load config: %v\n", err)
	}

	config := &config.Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Fatalf("Unable to unmarshal config: %v\n", err)
	}
	return config
}

func main() {
	// Set up database connection
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres_dev_password@localhost:5432/smarthome")
	database, err := db.New(dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer database.Close()

	log.Println("Connected to database successfully")

	// Initialize temperature service
	temperatureAPIURL := getEnv("TEMPERATURE_API_URL", "http://temperature-api:8081")
	temperatureService := services.NewTemperatureService(temperatureAPIURL)
	log.Printf("Temperature service initialized with API URL: %s\n", temperatureAPIURL)

	// Initialize router
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	apiRoutes := router.Group("/api/v1")

	var (
		secrets = mustLoadSecrets("secrets.yaml")
		config  = mustLoadConfig("config.yaml")
	)

	userAPI, err := userapi.NewClient(config.UserAPI.URL, secrets)
	if err != nil {
		log.Fatalf("Unable to create user API client: %v\n", err)
	}

	rabbitmqBroker, err := rabbitmq.NewRabbitMQBrokerController(&secrets.RabbitMQ)
	if err != nil {
		log.Fatalf("Unable to create rabbitmq broker: %v\n", err)
	}
	defer rabbitmqBroker.Close()

	broker, err := async.NewUserController(rabbitmqBroker)
	if err != nil {
		log.Fatalf("Unable to create user controller: %v\n", err)
	}

	// Register sensor routes
	sensorHandler := handlers.NewSensorHandler(database, temperatureService, broker, secrets, userAPI)
	sensorHandler.RegisterRoutes(apiRoutes)

	// Start server
	srv := &http.Server{
		Addr:    getEnv("PORT", ":8080"),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v\n", err)
	}

	log.Println("Server exited properly")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
