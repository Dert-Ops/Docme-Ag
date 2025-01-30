package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/Dert-Ops/Docme-Ag/internal/git"
)

// Commit işlemini başlatan fonksiyon
func RunCommitAgent() {
	reader := bufio.NewReader(os.Stdin)

	hasChanges, err := git.CheckGitStatus()
	if err != nil {
		fmt.Println("❌ Error checking git status:", err)
		return
	}
	if !hasChanges {
		fmt.Println("✅ No changes detected.")
		return
	}

	// **Yalnızca değişen satırları al**
	gitDiff, err := git.GetGitDiff()
	if err != nil {
		fmt.Println("❌ Error getting Git diff:", err)
		return
	}

	// AI tarafından üretilen commit mesajı almak için döngü
	var commitMessage string
	for {
		fmt.Println("🤖 Generating commit message using AI...")
		prompt := fmt.Sprintf("Analyze the following Git diff and suggest a Conventional Commit message:\n\n%s", gitDiff)
		commitMessage, err = gemini.GetGeminiResponse(prompt)
		if err != nil {
			fmt.Println("❌ Error getting AI commit message:", err)
			return
		}

		fmt.Println("\n📜 AI Suggested Commit Message:")
		fmt.Println(commitMessage)
		fmt.Println("\nDo you want to commit this change? (y/n/r)")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "y" || input == "Y" {
			break
		} else if input == "r" || input == "R" {
			fmt.Println("\n🔄 Regenerating commit message...")
			continue
		} else {
			fmt.Println("❌ Commit canceled.")
			return
		}
	}

	// Kullanıcı commit mesajını onayladıysa commit işlemini yap
	fmt.Println("✅ Committing changes...")
	err = git.CommitChanges(commitMessage)
	if err != nil {
		fmt.Println("❌ Error committing changes:", err)
		return
	}

	// Kullanıcıdan push için onay al
	fmt.Println("\n🚀 Do you want to push this commit to the repository? (y/n)")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

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
