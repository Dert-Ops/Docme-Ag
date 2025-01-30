package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Dert-Ops/Docme-Ag/internal/gemini"
	"github.com/fsnotify/fsnotify"
)

// İzlenecek dosya türleri
var fileExtensions = []string{".go", ".py", ".js", ".ts", ".java", ".cpp"}

// Kod dosyalarını izleyen agent
func WatchFiles(directory string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("Error creating watcher:", err)
	}
	defer watcher.Close()

	// Klasörü izlemeye al
	err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path) // Klasörleri izlemeye ekle
		}
		return nil
	})
	if err != nil {
		log.Fatal("Error walking directory:", err)
	}

	fmt.Println("🚀 Watching for file changes in:", directory)

	// Sürekli dosya değişikliklerini dinle
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Sadece belirli uzantılardaki dosyaları işle
			if isCodeFile(event.Name) && event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				fmt.Println("📂 File changed:", event.Name)
				processFile(event.Name)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}

// Belirtilen dosyanın uzantısını kontrol et
func isCodeFile(filename string) bool {
	for _, ext := range fileExtensions {
		if filepath.Ext(filename) == ext {
			return true
		}
	}
	return false
}

// Değişen dosyayı alıp Gemini'ye gönder
func processFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Println("Error reading file:", err)
		return
	}
	// Gemini API'ye kodu gönder ve öneri al
	fmt.Println("🤖 Sending file content to Gemini for suggestions...")
	response, err := gemini.GetGeminiResponse("Analyze this code and suggest improvements:\n\n" + string(content))
	if err != nil {
		log.Println("Error from Gemini API:", err)
		return
	}

	// Önerileri ekrana yazdır
	fmt.Println("✨ AI Suggestions for", filename, ":\n", response)
}
