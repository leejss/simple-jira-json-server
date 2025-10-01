package config

import "os"

type Config struct {
	Port           string
	DataPath       string
	AllowedOrigins []string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8080"),
		DataPath:       getEnv("DATA_PATH", "./data.json"),
		AllowedOrigins: []string{"*"},
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
