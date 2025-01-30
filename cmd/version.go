package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/Dert-Ops/Docme-Ag/internal/git"
)

// Semantic Versioning formatını kontrol eden regex
var semVerRegex = regexp.MustCompile(`\b\d+\.\d+\.\d+\b`)

// Versiyonlama işlemini yöneten fonksiyon
func RunVersioningAgent() {
	fmt.Println("🤖 Generating version number using AI...")

	// Tüm dosyaları oku
	allFilesContent := collectProjectFiles(".")
	reader := bufio.NewReader(os.Stdin)

	var newVersion string
	for {
		// Gemini API'den versiyon önerisi al
		prompt := fmt.Sprintf("Analyze this project and suggest a new Semantic Versioning number based on changes:\n\n%s", allFilesContent)
		response, err := gemini.GetGeminiResponse(prompt)
		if err != nil {
			fmt.Println("❌ Error getting AI versioning suggestion:", err)
			return
		}

		// Yanıtı temizle ve SemVer formatına uygun olup olmadığını kontrol et
		newVersion = ExtractVersionFromResponse(response)
		if newVersion == "" {
			fmt.Println("❌ AI did not provide a valid version number.")
			return
		}

		// Kullanıcıya önerilen versiyonu göster ve onay al
		fmt.Printf("\n📜 AI Suggested Version: v%s\n", newVersion)
		fmt.Println("\nDo you want to tag this version? (y/n/r)")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "y" || input == "Y" {
			break // Onaylandıysa döngüyü kır ve versiyon oluştur
		} else if input == "r" || input == "R" {
			fmt.Println("\n🔄 Regenerating version suggestion...")
			prompt = fmt.Sprintf(
				"The following version suggestion was incorrect. Generate a better Semantic Version number:\n\nPrevious version: v%s\n\nChanges:\n%s",
				newVersion, allFilesContent,
			)
			continue // Yeni versiyon önerisi al
		} else {
			fmt.Println("❌ Versioning canceled.")
			return
		}
	}

	// Kullanıcı versiyonu onayladıysa tag oluştur ve push et
	fmt.Printf("✅ Creating version tag v%s...\n", newVersion)
	err := git.CreateVersionTag(newVersion)
	if err != nil {
		fmt.Println("❌ Error creating version tag:", err)
		return
	}

	// Kullanıcıdan push için onay al
	fmt.Println("\n🚀 Do you want to push this tag to the repository? (y/n)")
	input, _ := reader.ReadString('\n')
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
}

// **Gemini yanıtından versiyon numarasını çıkart**
func ExtractVersionFromResponse(response string) string {
	matches := semVerRegex.FindStringSubmatch(response)
	if len(matches) > 0 {
		return matches[0] // İlk eşleşen versiyon numarasını döndür
	}
	return ""
}
