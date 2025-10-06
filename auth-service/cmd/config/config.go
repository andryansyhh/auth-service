package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string        `json:"APP_PORT"`
	DBDSN         string        `json:"DB_DSN"`
	JWTSecret     string        `json:"JWT_SECRET"`
	JWTExpiration time.Duration `json:"JWT_EXPIRATION"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system env")
	}

	exp, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		log.Fatal("invalid JWT_EXPIRATION format")
	}

	cfg := &Config{
		AppPort:       getEnv("APP_PORT", "8080"),
		DBDSN:         getEnv("DB_DSN", ""),
		JWTSecret:     getEnv("JWT_SECRET", "changeme"),
		JWTExpiration: exp,
	}

	if cfg.DBDSN == "" {
		log.Fatal("DB_DSN not found — please check your .env file")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
