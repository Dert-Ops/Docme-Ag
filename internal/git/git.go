package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// **Mevcut branch adını al**
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("❌ Error getting current branch: %v", err)
	}

	branch := strings.TrimSpace(out.String())
	if branch == "" {
		return "", fmt.Errorf("❌ No active Git branch found")
	}

	return branch, nil
}

// **Git Diff ile Değişiklikleri Al**
func GetGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--unified=0") // Sadece değişen satırları al
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("❌ Error getting Git diff: %v", err)
	}

	return out.String(), nil
}

// CreateVersionTag creates or updates a version tag
func CreateVersionTag(version string) error {
	// Check existing tags
	existingTagsCmd := exec.Command("git", "tag", "--list", "v"+version)
	var existingTags bytes.Buffer
	existingTagsCmd.Stdout = &existingTags
	err := existingTagsCmd.Run()
	if err != nil {
		return fmt.Errorf("❌ Error checking existing Git tags: %v", err)
	}

	// If tag exists, delete it
	if strings.TrimSpace(existingTags.String()) == "v"+version {
		fmt.Printf("⚠️  Warning: Tag v%s already exists. Deleting and recreating...\n", version)
		deleteCmd := exec.Command("git", "tag", "-d", "v"+version)
		if err := deleteCmd.Run(); err != nil {
			return fmt.Errorf("❌ Error deleting existing tag: %v", err)
		}
	}

	// Create new tag
	cmd := exec.Command("git", "tag", "-a", "v"+version, "-m", "Version "+version)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("❌ Error creating Git tag: %v", err)
	}

	fmt.Println("✅ Created new Git tag:", version)
	return nil
}

// Versiyon tag'ını remote repoya push et
func PushVersionTag(version string) error {
	cmd := exec.Command("git", "push", "origin", "v"+version)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("❌ Error pushing Git tag: %v", err)
	}

	fmt.Println("✅ Pushed Git tag:", version)
	return nil
}

// Git durumunu kontrol et
func CheckGitStatus() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("error checking git status: %v", err)
	}

	if len(output) == 0 {
		return false, nil
	}
	return true, nil
}

// **Git Push İşlemini Mevcut Branch İçin Yap**
func PushChanges() error {
	branch, err := GetCurrentBranch()
	if err != nil {
		return err // Mevcut branch alınamazsa hata döndür
	}

	fmt.Printf("📤 Pushing changes to remote repository (branch: %s)...\n", branch)
	cmd := exec.Command("git", "push", "origin", branch)
	cmd.Stdout = nil // Terminal çıktısını göster
	cmd.Stderr = nil

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("❌ Error pushing changes to branch %s: %v", branch, err)
	}

	fmt.Println("✅ Changes successfully pushed!")
	return nil
}
