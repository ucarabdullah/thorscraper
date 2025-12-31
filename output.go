package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// SavedData kaydedilen veri yapısı
type SavedData struct {
	URL         string    `json:"url"`
	StatusCode  int       `json:"status_code"`
	HTMLContent string    `json:"html_content"`
	ContentSize int       `json:"content_size"`
	Timestamp   time.Time `json:"timestamp"`
}

// SaveToJSON tüm sonuçları JSON dosyasına kaydeder
func SaveToJSON(results []ScanResult, filename string) error {
	outputDir := "logs"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	jsonFile := filepath.Join(outputDir, filename)

	// JSON formatına çevir
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return err
	}

	// JSON dosyasını kaydet
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

// sanitizeFilename URL'den güvenli dosya adı oluşturur
func sanitizeFilename(url string) string {
	// http:// veya https:// kaldır
	filename := strings.TrimPrefix(url, "http://")
	filename = strings.TrimPrefix(filename, "https://")

	// .onion uzantısını kaldır
	filename = strings.TrimSuffix(filename, ".onion")
	filename = strings.TrimSuffix(filename, ".onion/")

	// Geçersiz karakterleri temizle
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, ":", "_")
	filename = strings.ReplaceAll(filename, "*", "_")
	filename = strings.ReplaceAll(filename, "?", "_")
	filename = strings.ReplaceAll(filename, "\"", "_")
	filename = strings.ReplaceAll(filename, "<", "_")
	filename = strings.ReplaceAll(filename, ">", "_")
	filename = strings.ReplaceAll(filename, "|", "_")

	// Maksimum 50 karakter
	if len(filename) > 50 {
		filename = filename[:50]
	}

	return filename
}

// SaveHTMLToFile HTML içeriğini dosyaya kaydeder
func SaveHTMLToFile(name string, bodyBytes []byte) error {
	outputDir := "scraped_data"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("klasör oluşturulamadı: %v", err)
	}

	// İsimden güvenli dosya adı oluştur
	filename := sanitizeFilename(name)
	timestamp := time.Now().Format("20060102_150405")
	htmlFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.html", filename, timestamp))

	// HTML dosyasını kaydet
	if err := os.WriteFile(htmlFile, bodyBytes, 0644); err != nil {
		return fmt.Errorf("HTML kaydedilemedi: %v", err)
	}

	return nil
}
