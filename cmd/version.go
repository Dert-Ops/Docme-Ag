package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/Dert-Ops/Docme-Ag/internal/git"
)

// Semantic Versioning formatını kontrol eden regex
var semVerRegex = regexp.MustCompile(`\d+\.\d+\.\d+`)

// Versiyonlama işlemini yöneten fonksiyon
func RunVersioningAgent() {
	fmt.Println("🤖 Generating version number using AI...")

	// Gemini API'ye istek gönder
	response, err := gemini.GetGeminiResponse("Suggest a new Semantic Version number based on recent changes.")
	if err != nil {
		fmt.Println("❌ Error getting AI versioning suggestion:", err)
		return
	}

	// Yanıtı temizle ve SemVer formatına uygun olup olmadığını kontrol et
	newVersion := ExtractVersionFromResponse(response)
	if newVersion == "" {
		fmt.Println("❌ AI did not provide a valid version number.")
		return
	}

	// Kullanıcıdan onay al
	fmt.Printf("\n📜 AI Suggested Version: v%s\n", newVersion)
	fmt.Println("\nDo you want to tag this version? (y/n)")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = input[:len(input)-1]

	if input != "y" && input != "Y" {
		fmt.Println("❌ Versioning canceled.")
		return
	}

	// Git tag oluştur ve push et
	fmt.Printf("✅ Creating version tag v%s...\n", newVersion)
	err = git.CreateVersionTag(newVersion)
	if err != nil {
		fmt.Println("❌ Error creating version tag:", err)
		return
	}

	// Kullanıcıdan push için onay al
	fmt.Println("\n🚀 Do you want to push this tag to the repository? (y/n)")
	input, _ = reader.ReadString('\n')
	input = input[:len(input)-1]

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
