package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type URLurl struct {
	Name string
	URL  string
}

func ReadTargets(filename string) ([]URLurl, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("dosya okunamadı: %v", err)
	}

	lines := strings.Split(string(data), "\n")
	var targets []URLurl

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.Contains(line, "|") {
			parts := strings.SplitN(line, "|", 2)
			name := strings.TrimSpace(parts[0])
			url := strings.TrimSpace(parts[1])
			targets = append(targets, URLurl{Name: name, URL: url})
		} else {
			url := strings.TrimSpace(line)
			name := extractNameFromURL(url)
			targets = append(targets, URLurl{Name: name, URL: url})
		}
	}

	return targets, nil
}

func extractNameFromURL(url string) string {
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	if idx := strings.Index(url, ".onion"); idx != -1 {
		url = url[:idx]
	}

	if len(url) > 15 {
		url = url[:15] + "..."
	}

	return url
}

func main() {

	if err := InitLogger(); err != nil {
		log.Fatal("Logger başlatılamadı:", err)
	}
	defer CloseLogger()

	PrintDivider()
	fmt.Println(" THOR SCRAPER - made by Abdullah")
	PrintDivider()

	PrintStep(1, 5, "Hedef URL'ler okunuyor...")
	urls, err := ReadTargets("targets.yaml")
	if err != nil {
		PrintError("Dosya okuma hatası: %v", err)
		log.Fatal(err)
	}
	PrintSuccess("Toplam %d URL yüklendi", len(urls))

	PrintStep(2, 5, "Tor ağına bağlanılıyor...")
	client, torPort, err := CreateTorClient()
	if err != nil {
		PrintError("Tor bağlantısı kurulamadı: %v", err)
		log.Fatal(err)
	}

	PrintStep(3, 5, "Tor IP doğrulanıyor...")
	testTorIP(client)

	for {
		PrintStep(4, 5, "URL seçimi...")
		selectedURLs, exit := ShowMenu(urls)

		if exit {
			PrintInfo("Programdan çıkılıyor...")
			break
		}

		PrintSuccess("%d URL seçildi", len(selectedURLs))

		PrintStep(5, 5, "Tarama başlatılıyor...")
		fmt.Println()

		results := ScanAllURLs(client, selectedURLs, torPort)

		PrintDivider()
		PrintSummary(results)
		PrintDivider()

		SaveToJSON(results, "scan_results.json")
		fmt.Println()
		PrintSuccess("Tarama tamamlandı!")
		fmt.Println()

		fmt.Print("Yeni tarama yapmak için Enter'a basın (çıkmak için 'q' yazın): ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" || input == "Q" {
			PrintInfo("Programdan çıkılıyor...")
			break
		}

		fmt.Println()
	}

	PrintDivider()
	PrintSuccess("Tüm işlemler tamamlandı!")
	PrintDivider()
}

func testTorIP(client *http.Client) {
	LogInfo("check.torproject.org sorgulanıyor...")

	resp, err := client.Get("https://check.torproject.org/api/ip")
	if err != nil {
		PrintError("IP kontrolü başarısız")
		LogError("IP kontrolü hatası: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	LogSuccess("Tor IP: %s", string(body))
	PrintSuccess("Tor IP doğrulandı")

	resp2, err := client.Get("https://check.torproject.org")
	if err == nil {
		defer resp2.Body.Close()
		LogSuccess("Tor ağı bağlantısı test edildi: OK")
	} else {
		LogError("Tor test hatası: %v", err)
	}
}

func ShowMenu(urls []URLurl) ([]URLurl, bool) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n========== URL SEÇIM MENÜSÜ ==========")
		fmt.Printf("Toplam %d site bulundu:\n\n", len(urls))

		for i, url := range urls {
			fmt.Printf("[%d] %s\n", i+1, url.Name)
		}
		fmt.Printf("[16] Hepsini Tara\n\n")
		fmt.Printf("[0] Çıkış\n")

		fmt.Println("=======================================")
		fmt.Print("\nSeçiminizi yapın: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			return nil, true
		}

		if input == "" || input == "16" {
			fmt.Println("\n✓ Tüm siteler taranacak!")
			return urls, false
		}

		selected := []URLurl{}
		choices := strings.Split(input, ",")

		for _, choice := range choices {
			choice = strings.TrimSpace(choice)
			index, err := strconv.Atoi(choice)

			if err != nil || index < 1 || index > len(urls) {
				fmt.Printf("⚠ Geçersiz seçim: %s (atlandı)\n", choice)
				continue
			}

			selected = append(selected, urls[index-1])
		}

		if len(selected) == 0 {
			fmt.Println("\n⚠ Geçerli seçim yapılmadı! Lütfen tekrar deneyin.")
			continue
		}

		fmt.Println("\n✓ Seçilen siteler:")
		for _, url := range selected {
			fmt.Printf("  - %s\n", url.Name)
		}

		return selected, false
	}
}
