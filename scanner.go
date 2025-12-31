package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ScanResult struct {
	Name           string    `json:"name"`
	URL            string    `json:"url"`
	StatusCode     int       `json:"status_code"`
	Success        bool      `json:"success"`
	Error          string    `json:"error,omitempty"`
	Timestamp      time.Time `json:"timestamp"`
	SavedFile      string    `json:"saved_file,omitempty"`
	ScreenshotFile string    `json:"screenshot_file,omitempty"`
}

func ScanURL(client *http.Client, url URLurl, torPort string) ScanResult {
	result := ScanResult{
		Name:      url.Name,
		URL:       url.URL,
		Timestamp: time.Now(),
		Success:   false,
	}

	fmt.Printf("🔍 Scanning %s: %s", url.Name, url.URL)

	LogInfo("========================================")
	LogInfo("TARAMA: %s", url.Name)
	LogInfo("URL: %s", url.URL)
	LogDebug("HTTP Request oluşturuluyor...")

	req, err := http.NewRequest("GET", url.URL, nil)
	if err != nil {
		result.Error = err.Error()
		LogError("HTTP Request hatası: %v", err)
		fmt.Println(" ✗ Başarısız")
		return result
	}

	LogDebug("User-Agent ayarlanıyor...")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("DNT", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Upgrade-Insecure-Requests", "1")

	LogInfo("İstek gönderiliyor...")
	resp, err := client.Do(req)
	if err != nil {
		result.Error = err.Error()
		LogError("Bağlantı hatası: %v", err)
		fmt.Println(" ✗ Timeout/Fail")
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Success = true

	LogDebug("Response alındı - Status: %d", resp.StatusCode)

	var reader io.Reader = resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		LogDebug("Gzip sıkıştırması tespit edildi, açılıyor...")
		gzipReader, err := gzip.NewReader(resp.Body)
		if err == nil {
			reader = gzipReader
			defer gzipReader.Close()
		} else {
			LogWarning("Gzip açma hatası: %v", err)
		}
	}
	bodyBytes, _ := io.ReadAll(reader)
	bodySize := len(bodyBytes)

	LogSuccess("Sayfa indirildi - Status: %d, Boyut: %d bytes", resp.StatusCode, bodySize)

	fmt.Printf(" ✓ OK (%d)\n", resp.StatusCode)

	LogDebug("HTML kaydediliyor...")
	if err := SaveHTMLToFile(url.Name, bodyBytes); err != nil {
		LogWarning("HTML kaydedilemedi: %v", err)
	} else {
		result.SavedFile = sanitizeFilename(url.Name) + "_" + time.Now().Format("20060102_150405") + ".html"
		LogSuccess("HTML kaydedildi: %s", result.SavedFile)
	}

	if result.Success && result.StatusCode == 200 {
		LogInfo("Screenshot alınıyor...")
		fmt.Print("📸 Screenshot alınıyor...")

		screenshotFilename := sanitizeFilename(url.Name) + "_" + time.Now().Format("20060102_150405") + ".png"
		port, _ := strconv.Atoi(torPort)

		if TakeScreenshotSimple(url.URL, url.Name, port) {
			result.ScreenshotFile = screenshotFilename
			LogSuccess("Screenshot kaydedildi: %s", screenshotFilename)
			fmt.Println("Success ✓")
		} else {
			LogWarning("Screenshot başarısız")
			fmt.Println(" Error ✗")
		}
	}

	return result
}

func ScanAllURLs(client *http.Client, urls []URLurl, torPort string) []ScanResult {
	results := []ScanResult{}

	LogInfo("========================================")
	LogInfo("TOPLU TARAMA BAŞLIYOR - %d site", len(urls))
	LogInfo("========================================")

	for i, url := range urls {
		fmt.Printf("[%d/%d] ", i+1, len(urls))

		result := ScanURL(client, url, torPort)
		results = append(results, result)

		LogDebug("Rate limiting - 3 saniye bekleme")
		time.Sleep(3 * time.Second)
	}

	LogInfo("========================================")
	LogSuccess("TARAMA TAMAMLANDI")
	LogInfo("========================================")

	return results
}

func PrintSummary(results []ScanResult) {
	successCount := 0
	failCount := 0

	fmt.Println("\n📊 TARAMA SONUÇLARI")
	fmt.Println()

	fmt.Println("✓ Başarılı:")
	for _, result := range results {
		if result.Success {
			successCount++
			fmt.Printf("  • %s (Status: %d)\n", result.Name, result.StatusCode)
		}
	}

	for _, result := range results {
		if !result.Success {
			failCount++
		}
	}

	if failCount > 0 {
		fmt.Println("\n✗ Başarısız:")
		for _, result := range results {
			if !result.Success {
				fmt.Printf("  • %s\n", result.Name)
			}
		}
	}

	LogInfo("========== ÖZET ==========")
	LogInfo("Başarılı: %d", successCount)
	LogInfo("Başarısız: %d", failCount)
	LogInfo("Başarı oranı: %.2f%%", float64(successCount)/float64(len(results))*100)
}
