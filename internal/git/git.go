package git

import (
	"bytes"
	"fmt"
	"os/exec"
)

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

func CreateVersionTag(version string) error {
	cmd := exec.Command("git", "tag", "-a", "v"+version, "-m", "Version "+version)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error creating Git tag: %v", err)
	}

	fmt.Println("✅ Version", version, "created successfully!")
	return nil
}

// Yeni versiyon etiketi push et
func PushVersionTag(version string) error {
	cmd := exec.Command("git", "push", "origin", "v"+version)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error pushing Git tag: %v", err)
	}

	fmt.Println("✅ Version", version, "pushed successfully!")
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

// Git commit işlemi yap
func CommitChanges(commitMessage string) error {
	cmd := exec.Command("git", "add", ".")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding files: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", commitMessage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}

	fmt.Println("✅ Commit successful:", commitMessage)
	return nil
}

// Git push işlemi yap
func PushChanges() error {
	cmd := exec.Command("git", "push", "origin", "main")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error pushing changes: %v", err)
	}

	fmt.Println("✅ Push successful!")
	return nil
}
