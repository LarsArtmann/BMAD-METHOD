package config

import (
	"os"
	"strconv"
)

// Config holds the application configuration
type Config struct {
	Port    int    `json:"port" yaml:"port"`
	Version string `json:"version" yaml:"version"`
	Name    string `json:"name" yaml:"name"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Port:    8080,
		Version: "1.0.0",
		Name:    "demo-enterprise",
	}

	// Override with environment variables
	if port := os.Getenv("PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			cfg.Port = p
		}
	}

	if version := os.Getenv("VERSION"); version != "" {
		cfg.Version = version
	}

	return cfg, nil
}
