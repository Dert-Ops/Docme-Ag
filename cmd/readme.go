package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
)

// **README.md dosyasını güncelleyen fonksiyon**
func UpdateReadme(version, reason, summary string) error {
	readmePath := filepath.Join(".", "README.md")

	// **Mevcut README içeriğini oku**
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Println("❌ Error reading README.md:", err)
		return err
	}

	// **AI'ye güncellenmiş README.md oluşturması için prompt hazırla**
	context := "You are an AI assistant that updates README.md files with new version details."
	prompt := fmt.Sprintf(`
The following is the current README.md content:

%s

A new version has been released.

## 📌 New Version: v%s

- **Reason for Version Change:** %s
- **Summary of Changes:** %s

Update this README.md file to reflect the new version details in a structured and clear way.
Make sure to:
1. Clearly indicate that the new version is v%s.
2. Include a short description of the key changes.
3. Keep existing important content intact.
`, string(readmeContent), version, reason, summary, version)

	// **AI'den güncellenmiş README.md içeriğini al**
	updatedReadme, err := gemini.GetGeminiResponse(context, prompt)
	if err != nil {
		fmt.Println("❌ Error getting AI-suggested README.md:", err)
		return err
	}

	// **DEBUG: AI'den dönen içeriği yazdır**
	fmt.Println("\n🔍 AI Suggested README.md Content:")
	fmt.Println(updatedReadme)

	// **Eğer AI boş veya geçersiz bir içerik döndürdüyse hata ver**
	if updatedReadme == "" {
		fmt.Println("❌ AI response for README.md update was empty.")
		return fmt.Errorf("AI did not generate a valid README update")
	}

	// **README.md dosyasını güncelle**
	err = os.WriteFile(readmePath, []byte(updatedReadme), 0644)
	if err != nil {
		fmt.Println("❌ Error writing to README.md:", err)
		return err
	}

	// **Git commit işlemi**
	cmd := exec.Command("git", "add", "README.md")
	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Error adding README.md to Git:", err)
		return err
	}

	commitMessage := fmt.Sprintf(`docs: update README for v%s

- Updated version information in README.md
- Included explanation for the version bump
- Summarized key changes in the project
`, version)

	cmd = exec.Command("git", "commit", "-m", commitMessage)
	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Error committing README.md:", err)
		return err
	}

	fmt.Println("✅ README.md updated successfully!")
	return nil
}
