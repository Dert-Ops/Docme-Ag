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
var (
	SemVerRegex  = regexp.MustCompile(`(?i)VERSION:\s*v?(\d+\.\d+\.\d+)`)
	ReasonRegex  = regexp.MustCompile(`(?i)EXPLANATION:\s*([\s\S]+?)\n\nSUMMARY OF CHANGES:`)
	SummaryRegex = regexp.MustCompile(`(?i)SUMMARY OF CHANGES:\s*([\s\S]+)`)
)

// **Git üzerinden en son versiyon numarasını al**
func GetCurrentVersion() string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		fmt.Println("⚠️  Warning: No Git version tags found. Defaulting to v0.1.0")
		return "0.1.0"
	}

	return strings.TrimSpace(out.String()) // Versiyonu temizleyip döndür
}

// **Versiyonlama işlemini yöneten fonksiyon**
func RunVersioningAgent() {
	fmt.Println("🤖 Generating version number using AI...")

	// En son versiyon numarasını al
	currentVersion := GetCurrentVersion()

	// Tüm dosyalardaki değişiklikleri oku
	gitDiff, err := git.GetGitDiff()
	if err != nil {
		fmt.Println("❌ Error getting Git diff:", err)
		return
	}

	// **AI'ye yeni versiyon önerisi için prompt hazırla**
	context := "You are an AI assistant following Semantic Versioning principles."
	prompt := fmt.Sprintf(`
## Current Version: %s
## Changes:
%s

Analyze these changes and suggest a new Semantic Version number. 
Format: 
VERSION: X.Y.Z
EXPLANATION: 
- Reason 1
- Reason 2
- Reason 3
SUMMARY OF CHANGES:
- Change 1
- Change 2
- Change 3
`, currentVersion, gitDiff)

	aiResponse, err := gemini.GetGeminiResponse(context, prompt)

	// AI'den yeni versiyon önerisini al
	if err != nil {
		fmt.Println("❌ Error getting AI versioning suggestion:", err)
		return
	}

	fmt.Println("\n🔍 AI Response:")
	fmt.Println(aiResponse)

	// **AI yanıtını parse et ve versiyon numarası ile nedenini ayır**
	newVersion, reason, summary := ExtractVersionAndReason(aiResponse)
	if newVersion == "" {
		fmt.Println("❌ AI did not provide a valid version number.")
		return
	}

	// **Kullanıcıya önerilen versiyonu ve nedenini göster ve onay al**
	fmt.Printf("\n📜 AI Suggested Version: v%s\n", newVersion)
	fmt.Println("📝 Reason:", reason)
	fmt.Println("🔹 Summary:", summary)
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

		// **README.md dosyasını AI ile güncelle**
		err = UpdateReadme(fmt.Sprintf("New version released: v%s", newVersion), reason, summary)
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
func ExtractVersionAndReason(response string) (string, string, string) {
	versionMatch := SemVerRegex.FindStringSubmatch(response)
	reasonMatch := ReasonRegex.FindStringSubmatch(response)
	summaryMatch := SummaryRegex.FindStringSubmatch(response)

	var version, reason, summary string

	// **Versiyon Numarasını Al**
	if len(versionMatch) > 1 {
		version = strings.TrimSpace(versionMatch[1])
	}

	// **Açıklamayı Al**
	if len(reasonMatch) > 1 {
		reason = strings.TrimSpace(reasonMatch[1])
	}

	// **Özet Bilgisini Al**
	if len(summaryMatch) > 1 {
		summary = strings.TrimSpace(summaryMatch[1])
	}

	// **Eğer versiyon bulunamadıysa hata döndür**
	if version == "" {
		fmt.Println("❌ AI response did not contain a valid version number.")
		return "", "", ""
	}

	return version, reason, summary
}
