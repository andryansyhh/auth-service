package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	DBDSN         string
	JWTSecret     string
	JWTExpiration time.Duration
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	// Parsing durasi token JWT
	exp, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		log.Fatalf("Invalid format for JWT_EXPIRATION: %v", err)
	}

	cfg := &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		JWTSecret:     getEnv("JWT_SECRET", "defaultsecret"),
		DBDSN:         getEnv("DB_DSN", ""),
		JWTExpiration: exp,
	}

	if cfg.DBDSN == "" {
		log.Fatal("DB_DSN environment variable is required. Please check your .env file.")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
