package gemini

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Dert-Ops/Docme-Ag/config"
)

// Gemini 1.5 API URL
const apiURL = "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-pro:generateContent"

// API İstek Formatı
type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

// API Yanıt Formatı
type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// Gemini 1.5 API’ye mesaj gönder ve yanıt al
func GetGeminiResponse(prompt string) (string, error) {
	apiKey := config.GetAPIKey() // API anahtarını çevresel değişkenden al
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY is not set")
	}

	// Gemini API formatına uygun istek verisi
	requestBody, err := json.Marshal(GeminiRequest{
		Contents: []GeminiContent{
			{Parts: []GeminiPart{{Text: prompt}}},
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// HTTP isteğini oluştur
	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL+"?key="+apiKey, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Header bilgileri
	req.Header.Set("Content-Type", "application/json")

	// API isteğini gönder
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Yanıtı oku
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// JSON parse et
	var response GeminiResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %v", err)
	}

	// Yanıtı döndür
	if len(response.Candidates) == 0 || len(response.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("no response from Gemini API")
	}

	// **Yanıtı parse ederek anlamlı hale getir**
	return ParseGeminiResponse(response.Candidates[0].Content.Parts[0].Text), nil
}

// **Gemini Yanıtını Commit Mesajı veya Öneriler Haline Getir**
func ParseGeminiResponse(response string) string {
	lines := strings.Split(response, "\n")
	var formattedLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "*") || strings.HasPrefix(line, "-") {
			cleanedLine := strings.TrimPrefix(line, "* ")
			cleanedLine = strings.TrimPrefix(cleanedLine, "- ")
			formattedLines = append(formattedLines, "- "+cleanedLine) // Commit formatına uygun hale getir
		}
	}

	// Eğer hiçbir madde işaretli satır bulunamazsa, response'u direkt döndür
	if len(formattedLines) == 0 {
		return response
	}

	// **Commit mesajına uygun formatta string döndür**
	return fmt.Sprintf("refactor: AI-driven code improvements\n\n%s", strings.Join(formattedLines, "\n"))
}
