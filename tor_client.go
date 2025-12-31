package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

func CreateTorClient() (*http.Client, string, error) {
	ports := []string{"9050", "9150"}

	LogInfo("Tor servisi aranıyor...")

	for _, port := range ports {
		LogDebug("Port %s test ediliyor...", port)

		torProxy := "socks5://127.0.0.1:" + port

		proxyURL, err := url.Parse(torProxy)
		if err != nil {
			LogWarning("Port %s: URL parse hatası", port)
			continue
		}

		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			LogWarning("Port %s: Dialer hatası", port)
			continue
		}

		transport := &http.Transport{
			Dial:              dialer.Dial,
			DisableKeepAlives: true,
		}

		client := &http.Client{
			Transport: transport,
			Timeout:   time.Second * 40,
		}

		LogInfo("Port %s - Bağlantı test ediliyor...", port)
		resp, err := client.Get("https://check.torproject.org")
		if err == nil {
			resp.Body.Close()

			serviceName := "Tor Browser"
			if port == "9050" {
				serviceName = "Tor Standalone Service"
			}

			LogSuccess("Tor bulundu - Port: %s, Servis: %s", port, serviceName)
			LogSuccess("IP sızıntısı koruması: AKTİF")
			PrintSuccess("Tor bağlantısı başarılı! (Port: %s)", port)

			client.Timeout = time.Second * 120
			return client, port, nil
		}

		LogWarning("Port %s: Bağlantı başarısız", port)
	}

	LogError("Tor servisi bulunamadı (port 9050 ve 9150 denendi)")
	return nil, "", fmt.Errorf("Tor servisi bulunamadı")
}
