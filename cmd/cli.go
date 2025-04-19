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
		fmt.Println("  cm     - Commit changes using AI-generated messages")
		fmt.Println("  vs     - Generate new version using AI")
		fmt.Println("  readme - Update README.md with latest version details")
		return
	}

	switch os.Args[1] {
	case "cm":
		RunCommitAgent()
	case "vs":
		RunVersioningAgent()
	case "readme":
		if err := RunReadmeAgent(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	default:
		fmt.Println("Unknown command:", os.Args[1])
	}
}
