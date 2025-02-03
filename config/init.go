package config

import (
	"bufio"
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var (
	apiKey     string
	loadAPIKey sync.Once
)

// API anahtarını yükleyen fonksiyon
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
			apiKey = key

			// Key'i .env dosyasına kaydet
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

// API anahtarını döndüren fonksiyon
func GetAPIKey() string {
	loadKey()
	return apiKey
}
