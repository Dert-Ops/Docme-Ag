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
		// `.env` dosyasını farklı konumlarda arayacağız
		envPaths := []string{
			".env",                           // Proje kök dizininde
			os.Getenv("HOME") + "/.docm.env", // Kullanıcının home dizininde
			"/etc/docm.env",                  // Sistem genelinde
		}

		var loaded bool
		for _, path := range envPaths {
			if _, err := os.Stat(path); err == nil {
				err := godotenv.Load(path)
				if err == nil {
					fmt.Println("✅ Loaded environment variables from:", path)
					loaded = true
					break
				}
			}
		}

		if !loaded {
			fmt.Println("⚠️  Warning: No valid .env file found.")
		}

		// API anahtarını oku
		apiKey = os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			fmt.Println("⚠️  Warning: GEMINI_API_KEY is not set")
		}
	})
}

// API anahtarını döndüren fonksiyon
func GetAPIKey() string {
	loadKey()
	return apiKey
}
