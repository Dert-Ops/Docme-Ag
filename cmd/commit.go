package cmd

import (
	"fmt"
	"os/exec"

	"github.com/Dert-Ops/Docme-Ag/internal/chatgpt"
)

// Git commit işlemi yapan fonksiyon
func RunCommitAgent() {
	// Git değişiklikleri kontrol et
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error checking git status:", err)
		return
	}

	// Eğer değişiklik yoksa çık
	if len(output) == 0 {
		fmt.Println("No changes detected.")
		return
	}

	// ChatGPT'ye commit mesajı sorma
	fmt.Println("Generating commit message using AI...")
	commitMessage, err := chatgpt.GetChatGPTResponse("Suggest a Git commit message based on the latest code changes.")
	if err != nil {
		fmt.Println("Error getting AI commit message:", err)
		return
	}

	// Git commit işlemini gerçekleştir
	cmd = exec.Command("git", "add", ".")
	cmd.Run()

	cmd = exec.Command("git", "commit", "-m", commitMessage)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error committing changes:", err)
		return
	}

	fmt.Println("Commit successful:", commitMessage)
}
