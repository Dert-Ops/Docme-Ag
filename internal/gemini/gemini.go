package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/Dert-Ops/Docme-Ag/config"
)

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

// Gemini 1.5 API URL
const apiURL = "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-pro:generateContent"

// Semantic Versioning formatını kontrol eden regex
var commitRegex = regexp.MustCompile(`^(feat|fix|chore|docs|style|refactor|test|ci|build)(\(\w+\))?: (.+)$`)

// **Gemini Yanıtını Commit Mesajına Uygun Hale Getiren Fonksiyon**
func ParseCommitMessage(response string) string {
	lines := strings.Split(response, "\n")

	var commitHeader string
	var commitBody []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if commitRegex.MatchString(line) {
			commitHeader = line // İlk Conventional Commit formatına uyan satırı al
		} else if len(line) > 0 && len(commitHeader) > 0 {
			commitBody = append(commitBody, "- "+line) // Diğer açıklamaları madde işaretli hale getir
		}
	}

	// Eğer başlık bulunamazsa, varsayılan bir başlık belirle
	if commitHeader == "" {
		commitHeader = "chore: AI-generated commit message"
	}

	// Commit mesajını Conventional Commits formatına uygun hale getir
	return fmt.Sprintf("%s\n\n%s", commitHeader, strings.Join(commitBody, "\n"))
}

// Gemini API'ye mesaj gönder ve yanıtı Conventional Commits formatına uygun olarak parse et
func GetGeminiResponse(prompt string) (string, error) {
	apiKey := config.GetAPIKey()
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not set")
	}

	requestBody, err := json.Marshal(GeminiRequest{
		Contents: []GeminiContent{
			{Parts: []GeminiPart{{Text: prompt}}},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL+"?key="+apiKey, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var response GeminiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini API")
	}

	return response.Candidates[0].Content.Parts[0].Text, nil
}
