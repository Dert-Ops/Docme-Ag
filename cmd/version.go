package cmd

import (
	"fmt"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/Dert-Ops/Docme-Ag/internal/git"
)

// Versiyonlama işlemini yöneten fonksiyon
func RunVersioningAgent() {
	fmt.Println("Generating version number using AI...")

	// Gemini 1.5 API'den yeni Semantic Versioning numarası al
	newVersion, err := gemini.GetGeminiResponse("Suggest a new Semantic Version number based on recent changes.")
	if err != nil {
		fmt.Println("Error getting AI versioning suggestion:", err)
		return
	}

	// `git.go` içindeki fonksiyonu kullanarak versiyon oluştur
	err = git.CreateVersionTag(newVersion)
	if err != nil {
		fmt.Println("Error creating version tag:", err)
		return
	}
}
