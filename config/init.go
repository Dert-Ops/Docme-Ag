package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	apiKey     string
	loadAPIKey sync.Once
)

func loadKey() {
	loadAPIKey.Do(func() {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error: Environments can't get")
		}

		apiKey = os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			fmt.Println("Warning: GEMINI_API_KEY is not set")
		}
	})
}

// API anahtarını döndüren fonksiyon
func GetAPIKey() string {
	loadKey()
	return apiKey
}
