package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/Dert-Ops/Docme-Ag/config"
)

// **Gemini 1.5 API URL**
const apiURL = "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-pro:generateContent"

// **API İstek Formatı**
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Role  string       `json:"role"`
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

// **API Yanıt Formatı**
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Semantic Versioning formatını kontrol eden regex
var commitRegex = regexp.MustCompile(`^(feat|fix|chore|docs|style|refactor|test|ci|build)(\(\w+\))?: (.+)$`)

// **Gemini API'ye mesaj gönder ve yanıt al (`context` doğrudan `user` promptuna eklendi!)**
func GetGeminiResponse(context, prompt string) (string, error) {
	apiKey := config.GetAPIKey()
	if apiKey == "" {
		return "", fmt.Errorf("❌ ERROR: GEMINI_API_KEY is not set")
	}

	// **SYSTEM rolü kaldırıldı, context doğrudan prompt içinde**
	fullPrompt := fmt.Sprintf("%s\n\n%s", context, prompt)

	requestBody, err := json.Marshal(GeminiRequest{
		Contents: []GeminiContent{
			{
				Role: "user",
				Parts: []GeminiPart{
					{Text: fullPrompt},
				},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("❌ ERROR: Failed to marshal request body: %v", err)
	}

	// **HTTP isteğini oluştur**
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL+"?key="+apiKey, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("❌ ERROR: Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// **API'ye isteği gönder**
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("❌ ERROR: Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// **Yanıtı oku**
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("❌ ERROR: Failed to read response body: %v", err)
	}

	// **API yanıtını ekrana yazdır (debug için)**
	fmt.Println("\n🔍 Raw API Response:", string(body))

	// **Yanıtı JSON olarak çözümle**
	var response GeminiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("❌ ERROR: Failed to parse JSON response: %v", err)
	}

	// **Yanıt boşsa hata ver**
	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("❌ ERROR: No response from Gemini API")
	}

	// **Yanıtı döndür**
	return response.Candidates[0].Content.Parts[0].Text, nil
}
