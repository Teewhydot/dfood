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
}

type DatabaseConfig struct {
	Driver     string `yaml:"driver"`
	Datasource string `yaml:"datasource"`
}

func New() (*Config, error) {
	env := getEnvOrDefault("APP_ENV", "dev")
	configFile := fmt.Sprintf("config/config.%s.yaml", env)
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configFile, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
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
