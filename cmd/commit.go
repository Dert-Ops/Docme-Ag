package cmd

import (
	"fmt"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/Dert-Ops/Docme-Ag/internal/git"
)

// Git commit işlemi yapan fonksiyon
func RunCommitAgent() {
	// Git değişikliklerini kontrol et
	hasChanges, err := git.CheckGitStatus()
	if err != nil {
		fmt.Println("Error checking git status:", err)
		return
	}
	if !hasChanges {
		fmt.Println("No changes detected.")
		return
	}

	// Gemini 1.5 API'den commit mesajı al
	fmt.Println("Generating commit message using AI...")
	commitMessage, err := gemini.GetGeminiResponse("Suggest a Git commit message based on the latest code changes.")
	if err != nil {
		fmt.Println("Error getting AI commit message:", err)
		return
	}

	// `git.go` içindeki fonksiyonu kullanarak commit işlemi yap
	err = git.CommitChanges(commitMessage)
	if err != nil {
		fmt.Println("Error committing changes:", err)
		return
	}
}
