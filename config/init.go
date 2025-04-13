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

// **API Anahtarını Yükleme Fonksiyonu**
func loadKey() {
	loadAPIKey.Do(func() {
		// `.env` dosyasını yüklemeyi dene
		err := godotenv.Load()
		if err != nil {
			fmt.Println("⚠️ Warning: Could not load .env file. If this is your first time using Docme-Ag, set your API key.")
		}

		// Çevresel değişkenlerden API anahtarını al
		apiKey = os.Getenv("GEMINI_API_KEY")

		// **Eğer API Key Set Edilmemişse Kullanıcıdan İste**
		if strings.TrimSpace(apiKey) == "" {
			fmt.Println("\n🔑 GEMINI API Key is not set.")
			fmt.Print("👉 Please enter your GEMINI API Key: ")

			reader := bufio.NewReader(os.Stdin)
			key, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(key) // Boşlukları temizle

			// Kullanıcı boş API key girdiyse hata ver
			if apiKey == "" {
				fmt.Println("❌ Error: No API key provided. Exiting...")
				os.Exit(1)
			}

			// **Geçici olarak ENV değişkenine ata**
			os.Setenv("GEMINI_API_KEY", apiKey)

			// **API Key'i `.env` dosyasına kaydet**
			saveToEnvFile(apiKey)
		}
	})
}

// **ENV Dosyasına API Key Kaydetme**
func saveToEnvFile(key string) {
	file, err := os.OpenFile(".env", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("⚠️ Warning: Could not save API key to .env file.")
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("GEMINI_API_KEY=%s\n", key))
	if err != nil {
		fmt.Println("⚠️ Warning: Could not write API key to .env file.")
	} else {
		fmt.Println("✅ API Key saved successfully in .env file!")
	}
}

// **API Anahtarını Döndüren Fonksiyon**
func GetAPIKey() string {
	loadKey()
	return apiKey
}
