package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env string

	API APIConfig
	DB  DBConfig
}

type APIConfig struct {
	Port int
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() *Config {
	_ = godotenv.Load() // загружаем .env, если есть

	cfg := &Config{
		Env: getEnv("ENV", "local"),
		API: APIConfig{
			Port: getEnvAsInt("API_PORT", 8080),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "defi"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if v := os.Getenv(key); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid value for %s", key)
		}
		return i
	}
	return defaultValue
}
