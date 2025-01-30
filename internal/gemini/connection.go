package gemini

import (
	"context"
	"fmt"
	"log"

	"github.com/Dert-Ops/Docme-Ag/config"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}

func CreateGenerativeModel() error {

	fmt.Println(config.GetAPIKey())

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(config.GetAPIKey()))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(""))
	if err != nil {
		log.Fatal(err)
		return err
	}

	printResponse(resp)
	return nil
}
