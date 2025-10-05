package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AppName  string         `yaml:"app_name"`
	Env      string         `yaml:"env"`
	Port     int            `yaml:"port"`
	DB       DatabaseConfig `yaml:"db"`
	LogLevel string         `yaml:"log_level"`
	SendGrid SendGridConfig `yaml:"sendgrid"`
}

type DatabaseConfig struct {
	Driver     string `yaml:"driver"`
	Datasource string `yaml:"datasource"`
}

type SendGridConfig struct {
	APIKey    string `yaml:"api_key"`
	FromEmail string `yaml:"from_email"`
	FromName  string `yaml:"from_name"`
}

func New() (*Config, error) {
	env := getEnvOrDefault("APP_ENV", "dev")
	configFile := fmt.Sprintf("config/config.%s.yaml", env)
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	// Replace environment variables in config
	configContent := string(data)
	configContent = expandEnvVars(configContent)

	var cfg Config
	if err := yaml.Unmarshal([]byte(configContent), &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// expandEnvVars replaces ${VAR_NAME} with environment variable values
func expandEnvVars(content string) string {
	return os.Expand(content, func(key string) string {
		return os.Getenv(key)
	})
}
