package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AzureStorageAccount   string
	AzureStorageKey       string
	AzureStorageContainer string
	Port                  string
	JWTSecret             string
	DatabaseURL           string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	config := &Config{
		AzureStorageAccount:   getEnvOrDefault("AZURE_STORAGE_ACCOUNT", ""),
		AzureStorageKey:       getEnvOrDefault("AZURE_STORAGE_ACCESS_KEY", ""),
		AzureStorageContainer: getEnvOrDefault("AZURE_STORAGE_CONTAINER", ""),
		Port:                  getEnvOrDefault("PORT", "8080"),
		JWTSecret:             getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		DatabaseURL:           getEnvOrDefault("DATABASE_URL", "postgres://postgres:password@localhost:5432/cloudstorage?sslmode=disable"),
	}

	validateConfig(config)
	return config
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func validateConfig(config *Config) {
	if config.AzureStorageAccount == "" {
		log.Fatal("AZURE_STORAGE_ACCOUNT is required")
	}
	if config.AzureStorageKey == "" {
		log.Fatal("AZURE_STORAGE_ACCESS_KEY is required")
	}
	if config.AzureStorageContainer == "" {
		log.Fatal("AZURE_STORAGE_CONTAINER is required")
	}
}
