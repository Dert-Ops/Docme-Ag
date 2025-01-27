package git

import (
	"fmt"
	"os/exec"
)

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
	// Değişiklikleri ekle
	cmd := exec.Command("git", "add", ".")
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error adding files: %v", err)
	}

	// Commit işlemi yap
	cmd = exec.Command("git", "commit", "-m", commitMessage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error committing changes: %v", err)
	}

	fmt.Println("✅ Commit successful:", commitMessage)
	return nil
}

// Yeni versiyon etiketi oluştur ve push et
func CreateVersionTag(version string) error {
	cmd := exec.Command("git", "tag", "-a", "v"+version, "-m", "Version "+version)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error creating Git tag: %v", err)
	}

	cmd = exec.Command("git", "push", "origin", "v"+version)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error pushing Git tag: %v", err)
	}

	fmt.Println("✅ Version", version, "pushed successfully!")
	return nil
}
