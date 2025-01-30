package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/Dert-Ops/Docme-Ag/internal/git"
)

// Desteklenen dosya uzantıları
var supportedExtensions = map[string]struct{}{
	".go": {}, ".py": {}, ".js": {}, ".ts": {},
	".java": {}, ".cpp": {}, ".c": {}, ".cs": {},
}

// Commit işlemini başlatan fonksiyon
func RunCommitAgent() {
	hasChanges, err := git.CheckGitStatus()
	if err != nil {
		fmt.Println("❌ Error checking git status:", err)
		return
	}
	if !hasChanges {
		fmt.Println("✅ No changes detected.")
		return
	}
	reader := bufio.NewReader(os.Stdin)

	// Tüm proje dosyalarını oku
	allFilesContent := collectProjectFiles(".")

	// AI tarafından üretilen commit mesajı almak için döngü
	var commitMessage string
	for {
		fmt.Println("🤖 Generating commit message using AI...")
		prompt := fmt.Sprintf("Analyze these code changes and suggest a Conventional Commit message:\n\n%s", allFilesContent)
		commitMessage, err = gemini.GetGeminiResponse(prompt)
		if err != nil {
			fmt.Println("❌ Error getting AI commit message:", err)
			return
		}

		fmt.Println("\n📜 AI Suggested Commit Message:\n")
		fmt.Println(commitMessage)
		fmt.Println("\nDo you want to commit this change? (y/n/r)")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "y" || input == "Y" {
			break // Onaylandıysa döngüyü kır ve commit işlemi yap
		} else if input == "r" || input == "R" {
			fmt.Println("\n🔄 Regenerating commit message...")
			prompt = fmt.Sprintf(
				"The following commit message was not correct. Generate a better Conventional Commit message:\n\nPrevious commit message:\n%s\n\nChanges:\n%s",
				commitMessage, allFilesContent,
			)
			continue // Yeni commit mesajı al
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

// **Tüm proje dosyalarını oku ve içeriği tek bir string olarak döndür**
func collectProjectFiles(rootDir string) string {
	var allFilesContent strings.Builder

	// Dosya ve dizinleri gez
	filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil // Klasörleri atla
		}
		if _, exists := supportedExtensions[filepath.Ext(path)]; !exists {
			return nil // Desteklenmeyen dosya türlerini atla
		}

		// Dosya içeriğini oku
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		// İçeriği ekle
		allFilesContent.WriteString(fmt.Sprintf("\n\nFile: %s\n%s", path, string(content)))
		return nil
	})

	return allFilesContent.String()
}
