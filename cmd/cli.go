package cmd

import (
	"fmt"
	"os"
)

// CLI giriş noktası
func ExecuteCommand() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: docm <command>")
		fmt.Println("Commands:")
		fmt.Println("  cm   - Commit changes using AI-generated messages")
		fmt.Println("  vs   - Generate new version using AI")
		return
	}

	switch os.Args[1] {
	case "cm":
		RunCommitAgent()
	case "vs":
		RunVersioningAgent()
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
