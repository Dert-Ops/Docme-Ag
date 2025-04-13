package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
)

// README.md'yi güncelleyen fonksiyon
func UpdateReadme(commitMessage, reason, summary string) error {
	readmePath := "README.md"

	// README.md mevcut mu?
	_, err := os.Stat(readmePath)
	if os.IsNotExist(err) {
		fmt.Println("⚠️ README.md not found. Creating a new one...")
		file, err := os.Create(readmePath)
		if err != nil {
			fmt.Println("❌ Error creating README.md:", err)
			return err
		}
		file.Close()
	}

	// README.md içeriğini oku
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Println("❌ Error reading README.md:", err)
		return err
	}

	// **AI'ye yeni versiyon bilgilerini sorarken context ekle**
	context := "You are an AI assistant responsible for updating a project's README.md file. Ensure clarity and structure while incorporating the latest version details."

	currentVersion := GetCurrentVersion()

	prompt := fmt.Sprintf(`
The following is the current README.md content:

%s

A new version has been released.

## New Version Details
- **Version:** %s
- **Commit Message:** %s
- **Reason for Version Change:** %s
- **Summary of Changes:** %s

Update this README.md file to reflect the new version details in a structured and clear way.
`, string(readmeContent), currentVersion, commitMessage, reason, summary)

	// AI'den güncellenmiş README içeriğini al
	updatedReadme, err := gemini.GetGeminiResponse(context, prompt)
	if err != nil {
		fmt.Println("❌ Error getting AI-suggested README.md:", err)
		return err
	}

	// Eğer AI boş veya geçersiz veri döndürdüyse, eski halini koru
	if updatedReadme == "" || len(updatedReadme) < 50 {
		fmt.Println("⚠️ AI did not return a valid README update. Keeping the old version.")
		return nil
	}

	// Eğer yeni içerik aynıysa, güncelleme yapma
	if string(readmeContent) == updatedReadme {
		fmt.Println("✅ README.md is already up-to-date. No changes made.")
	} else {
		// README.md'yi güncelle
		err = os.WriteFile(readmePath, []byte(updatedReadme), 0644)
		if err != nil {
			fmt.Println("❌ Error writing to README.md:", err)
			return err
		}

		// Git commit işlemi
		cmd := exec.Command("git", "add", "README.md")
		err = cmd.Run()
		if err != nil {
			fmt.Println("❌ Error adding README.md to Git:", err)
			return err
		}

		cmd = exec.Command("git", "commit", "-m", "docs: update README with new version details")
		err = cmd.Run()
		if err != nil {
			fmt.Println("❌ Error committing README.md:", err)
			return err
		}

		fmt.Println("✅ README.md updated successfully!")
	}
	return nil
}

// this func behave like router

// func RunReadmeAgent() error {

// 	// way 1: If there is any readme in the project

// 	// way 2: If there is already readme file

// 	return nil
// }
