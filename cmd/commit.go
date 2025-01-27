package cmd

import (
	"fmt"

	"github.com/Dert-Ops/Docme-Ag/internal/chatgpt"
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

	// ChatGPT'den commit mesajı al
	fmt.Println("Generating commit message using AI...")
	commitMessage, err := chatgpt.GetChatGPTResponse("Suggest a Git commit message based on the latest code changes.")
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
