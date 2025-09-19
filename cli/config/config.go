package config

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

type Config struct {
	JiraApiToken string
	JiraBaseURL  string
	Username     string
	RawOutputDir string
	OutputDir    string
}

func LoadConfig() (*Config, error) {

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		JiraApiToken: getRequiredEnv("JIRA_API_TOKEN"),
		JiraBaseURL:  getRequiredEnv("JIRA_BASE_URL"),
		Username:     getRequiredEnv("JIRA_USERNAME"),
		RawOutputDir: getOptionalEnv("JIRA_RAW_OUTPUT_DIR", "output/raw"),
		OutputDir:    getOptionalEnv("JIRA_OUTPUT_DIR", "output/data"),
	}, nil
}

func getOptionalEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getRequiredEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("required env var %s is not set", key)
	}
	return value
}
