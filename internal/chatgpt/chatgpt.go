package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiURL = "https://api.openai.com/v1/chat/completions"

// API isteği için veri yapısı
type OpenAIRequest struct {
	Model    string        `json:"model"`
	Messages []MessageItem `json:"messages"`
}

type MessageItem struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// API'den dönen yanıt
type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// ChatGPT'ye mesaj gönder ve yanıt al
func GetChatGPTResponse(prompt string) (string, error) {
	// API anahtarını çevresel değişkenden al
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("OPENAI_API_KEY is not set")
	}

	// API için mesaj formatı
	messages := []MessageItem{
		{"system", "You are an AI assistant helping with Git commits and versioning."},
		{"user", prompt},
	}

	// JSON formatına çevir
	requestBody, err := json.Marshal(OpenAIRequest{
		Model:    "gpt-4",
		Messages: messages,
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// HTTP isteği oluştur
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Header bilgilerini ekle
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	// İsteği gönder
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Yanıtı oku (ioutil yerine io kullanıldı)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// JSON parse et
	var response OpenAIResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Yanıt kontrolü
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}
