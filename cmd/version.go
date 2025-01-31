package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/Dert-Ops/Docme-Ag/internal/git"
)

// Semantic Versioning formatını kontrol eden regex
var semVerRegex = regexp.MustCompile(`\b\d+\.\d+\.\d+\b`)

func GetCurrentVersion() string {
	// `git describe --tags` ile en son versiyon tag'ini al
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("⚠️  Warning: No Git version tags found. Defaulting to v0.1.0")
		return "0.1.0" // Eğer hiç tag yoksa varsayılan değer
	}

	return strings.TrimSpace(out.String()) // Versiyonu temizleyip döndür
}

// Versiyonlama işlemini yöneten fonksiyon
func RunVersioningAgent() {
	fmt.Println("🤖 Generating version number using AI...")

	// En son versiyon numarasını `git` üzerinden al
	currentVersion := GetCurrentVersion()

	// Tüm dosyalardaki değişiklikleri oku
	gitDiff, err := git.GetGitDiff()
	if err != nil {
		fmt.Println("❌ Error getting Git diff:", err)
		return
	}

	// AI'ye yeni versiyon önerisi için prompt hazırla
	prompt := fmt.Sprintf(
		"The current version is %s. Analyze the following Git diff and suggest a new Semantic Versioning number based on the changes:\n\n%s",
		currentVersion, gitDiff,
	)

	// AI'den yeni versiyon önerisini al
	newVersion, err := gemini.GetGeminiResponse(prompt)
	if err != nil {
		fmt.Println("❌ Error getting AI versioning suggestion:", err)
		return
	}

	// AI yanıtından versiyon numarasını temizle
	newVersion = ExtractVersionFromResponse(newVersion)
	if newVersion == "" {
		fmt.Println("❌ AI did not provide a valid version number.")
		return
	}

	// Kullanıcıya önerilen versiyonu göster ve onay al
	fmt.Printf("\n📜 AI Suggested Version: v%s\n", newVersion)
	fmt.Println("\nDo you want to tag this version? (y/n/r)")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "y" || input == "Y" {
		// Yeni versiyon için Git tag oluştur
		fmt.Printf("✅ Creating version tag v%s...\n", newVersion)
		err = git.CreateVersionTag(newVersion)
		if err != nil {
			fmt.Println("❌ Error creating version tag:", err)
			return
		}

		// 📜 **README.md dosyasını AI ile güncelle**
		err = UpdateReadme("New version released", newVersion)
		if err != nil {
			fmt.Println("❌ Error updating README.md:", err)
			return
		}

		// Kullanıcıdan push için onay al
		fmt.Println("\n🚀 Do you want to push this tag to the repository? (y/n)")
		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "y" || input == "Y" {
			fmt.Println("📤 Pushing version tag to remote repository...")
			err = git.PushVersionTag(newVersion)
			if err != nil {
				fmt.Println("❌ Error pushing version tag:", err)
				return
			}
			fmt.Println("✅ Version tag successfully pushed!")
		} else {
			fmt.Println("❌ Push canceled.")
		}
	} else if input == "r" || input == "R" {
		fmt.Println("\n🔄 Regenerating version suggestion...")
		RunVersioningAgent() // Yeniden başlat
	} else {
		fmt.Println("❌ Versioning canceled.")
	}
}

// **Gemini yanıtından versiyon numarasını çıkart**
func ExtractVersionFromResponse(response string) string {
	matches := semVerRegex.FindStringSubmatch(response)
	if len(matches) > 0 {
		return matches[0] // İlk eşleşen versiyon numarasını döndür
	}
	return ""
}
