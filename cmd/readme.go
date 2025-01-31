package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
)

func UpdateReadme(commitMessage, version string) error {
	readmePath := "README.md"

	// README.md dosyasını oku
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Println("❌ Error reading README.md:", err)
		return err
	}

	// Versiyon numarasını README.md'ye ekleyerek AI'ye prompt gönder
	prompt := fmt.Sprintf(`
The following is the current README.md content:

%s

The latest version is now v%s.

Update this README.md file to reflect the new version information in a structured and clear way.
`, string(readmeContent), version)

	// AI'den güncellenmiş README.md içeriğini al
	updatedReadme, err := gemini.GetGeminiResponse(prompt)
	if err != nil {
		fmt.Println("❌ Error getting AI-suggested README.md:", err)
		return err
	}

	// README.md dosyasını güncelle
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

	cmd = exec.Command("git", "commit", "-m", "docs: update README with new version information")
	err = cmd.Run()
	if err != nil {
		fmt.Println("❌ Error committing README.md:", err)
		return err
	}

	fmt.Println("✅ README.md updated successfully!")
	return nil
}
