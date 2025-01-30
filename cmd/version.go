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
// Versiyonlama işlemini yöneten fonksiyon
func RunVersioningAgent() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🤖 Generating version number using AI...")

	// **Yalnızca değişen satırları al**
	gitDiff, err := git.GetGitDiff()
	if err != nil {
		fmt.Println("❌ Error getting Git diff:", err)
		return
	}

	var newVersion string
	var prompt string
	for {
		// Gemini API'den versiyon önerisi al
		if prompt == "" {
			prompt = fmt.Sprintf("Analyze the following Git diff and suggest a new Semantic Versioning number:\n\n%s", gitDiff)
		}
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
			break
		} else if input == "r" || input == "R" {
			fmt.Println("\n🔄 Regenerating version suggestion...")
			prompt = fmt.Sprintf(
				"The following version suggestion was incorrect. Generate a better Semantic Version number:\n\nPrevious version: v%s\n\nChanges:\n%s",
				newVersion, gitDiff,
			)
			continue
		} else {
			fmt.Println("❌ Versioning canceled.")
			return
		}
	}

	// Kullanıcı versiyonu onayladıysa tag oluştur ve push et
	fmt.Printf("✅ Creating version tag v%s...\n", newVersion)
	err = git.CreateVersionTag(newVersion)
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
