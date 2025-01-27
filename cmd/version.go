package cmd

import (
	"fmt"
	"os/exec"

	"github.com/Dert-Ops/Docme-Ag/internal/chatgpt"
)

// Versiyonlama işlemini yöneten fonksiyon
func RunVersioningAgent() {
	fmt.Println("Generating version number using AI...")

	// OpenAI'ye versiyon sorma
	newVersion, err := chatgpt.GetChatGPTResponse("Suggest a new Semantic Version number based on recent changes.")
	if err != nil {
		fmt.Println("Error getting AI versioning suggestion:", err)
		return
	}

	// Yeni versiyonu Git Tag olarak ekle
	cmd := exec.Command("git", "tag", "-a", "v"+newVersion, "-m", "Version "+newVersion)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error creating Git tag:", err)
		return
	}

	cmd = exec.Command("git", "push", "origin", "v"+newVersion)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error pushing Git tag:", err)
		return
	}

	fmt.Println("Version", newVersion, "pushed successfully!")
}
