package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/Dert-Ops/Docme-Ag/internal/git"
)

// Git commit işlemi yapan fonksiyon (AI ile commit mesajı oluşturma)
func RunCommitAgent() {
	// Git değişikliklerini kontrol et
	hasChanges, err := git.CheckGitStatus()
	if err != nil {
		fmt.Println("❌ Error checking git status:", err)
		return
	}
	if !hasChanges {
		fmt.Println("✅ No changes detected.")
		return
	}

	// Gemini 1.5 API'den commit mesajı al
	fmt.Println("🤖 Generating commit message using AI...")
	commitMessage, err := gemini.GetGeminiResponse("Analyze the latest code changes and suggest a Git commit message.")
	if err != nil {
		fmt.Println("❌ Error getting AI commit message:", err)
		return
	}

	// Kullanıcıya commit mesajını göster ve onay al
	fmt.Println("\n📜 AI Suggested Commit Message:\n")
	fmt.Println(commitMessage)
	fmt.Println("\nDo you want to commit this change? (y/n)")

	// Kullanıcıdan giriş al
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1] // Yeni satır karakterini kaldır

	if input != "y" && input != "Y" {
		fmt.Println("❌ Commit canceled.")
		return
	}

	// `git.go` içindeki fonksiyonu kullanarak commit işlemi yap
	fmt.Println("✅ Committing changes...")
	err = git.CommitChanges(commitMessage)
	if err != nil {
		fmt.Println("❌ Error committing changes:", err)
		return
	}

	// Kullanıcıdan push için onay al
	fmt.Println("\n🚀 Do you want to push this commit to the repository? (y/n)")
	input, _ = reader.ReadString('\n')
	input = input[:len(input)-1]

	if input == "y" || input == "Y" {
		fmt.Println("📤 Pushing changes to remote repository...")
		err = git.PushChanges()
		if err != nil {
			fmt.Println("❌ Error pushing changes:", err)
			return
		}
		fmt.Println("✅ Changes successfully pushed!")
	} else {
		fmt.Println("❌ Push canceled.")
	}
}
