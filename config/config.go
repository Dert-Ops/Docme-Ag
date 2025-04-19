package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

var (
	apiKey     string
	loadAPIKey sync.Once
)

// loadKey loads the API key from .env file or prompts user for input
func loadKey() {
	loadAPIKey.Do(func() {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("⚠️ Warning: Could not load .env file. If this is your first time using Docme-Ag, set your API key.")
		}

		apiKey = os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			fmt.Println("\n🔑 GEMINI API Key is not set.")
			fmt.Print("👉 Please enter your GEMINI API Key: ")
			reader := bufio.NewReader(os.Stdin)
			key, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(key)

			// Save key to .env file
			file, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err == nil {
				_, _ = file.WriteString(fmt.Sprintf("GEMINI_API_KEY=%s\n", apiKey))
				file.Close()
				fmt.Println("✅ API Key saved successfully in .env file!")
			} else {
				fmt.Println("❌ Failed to save API Key. Please set it manually.")
			}
		}
	})
}

// GetAPIKey returns the Gemini API key from environment variables or prompts for input
func GetAPIKey() string {
	loadKey()
	return strings.TrimSpace(apiKey)
}

// SetAPIKey sets the Gemini API key as an environment variable and saves to .env
func SetAPIKey(key string) error {
	apiKey = strings.TrimSpace(key)
	err := os.Setenv("GEMINI_API_KEY", apiKey)
	if err != nil {
		return err
	}

	// Save to .env file
	file, err := os.OpenFile(".env", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("GEMINI_API_KEY=%s\n", apiKey))
	return err
} 