// Package config provides configuration management for the UniFi MCP server.
package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the UniFi controller configuration.
type Config struct {
	Host           string
	Username       string
	Password       string
	Port           int
	Site           string
	VerifySSL      bool
	ControllerType string // auto, proxy, direct
	LogLevel       string
}

// Load reads configuration from environment variables.
func Load() *Config {
	cfg := &Config{
		Host:           getEnv("UNIFI_HOST", ""),
		Username:       getEnv("UNIFI_USERNAME", ""),
		Password:       getEnv("UNIFI_PASSWORD", ""),
		Port:           getEnvInt("UNIFI_PORT", 443),
		Site:           getEnv("UNIFI_SITE", "default"),
		VerifySSL:      getEnvBool("UNIFI_VERIFY_SSL", false),
		ControllerType: getEnv("UNIFI_CONTROLLER_TYPE", "auto"),
		LogLevel:       getEnv("UNIFI_MCP_LOG_LEVEL", "INFO"),
	}
	return cfg
}

// Validate checks if required configuration is present.
func (c *Config) Validate() error {
	if c.Host == "" {
		return &ConfigError{Field: "UNIFI_HOST", Message: "is required"}
	}
	if c.Username == "" {
		return &ConfigError{Field: "UNIFI_USERNAME", Message: "is required"}
	}
	if c.Password == "" {
		return &ConfigError{Field: "UNIFI_PASSWORD", Message: "is required"}
	}
	return nil
}

// BaseURL returns the base URL for the UniFi controller.
func (c *Config) BaseURL() string {
	return "https://" + c.Host + ":" + strconv.Itoa(c.Port)
}

// ConfigError represents a configuration error.
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return e.Field + " " + e.Message
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		lower := strings.ToLower(value)
		return lower == "true" || lower == "1" || lower == "yes"
	}
	return defaultValue
}
