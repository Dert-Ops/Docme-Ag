package config

import (
	"os"
	"strings"
)

// GetAPIKey returns the Gemini API key from environment variables
func GetAPIKey() string {
	key := os.Getenv("GEMINI_API_KEY")
	return strings.TrimSpace(key)
}

// SetAPIKey sets the Gemini API key as an environment variable
func SetAPIKey(key string) error {
	return os.Setenv("GEMINI_API_KEY", key)
} 