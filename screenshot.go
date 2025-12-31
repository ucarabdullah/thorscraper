package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/chromedp/chromedp"
)

func TakeScreenshot(url, filename string, torPort int) error {
	screenshotDir := "screenshots"
	if err := os.MkdirAll(screenshotDir, 0755); err != nil {
		return fmt.Errorf("klasör oluşturulamadı: %v", err)
	}

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer(fmt.Sprintf("socks5://127.0.0.1:%d", torPort)),
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-images", false),
		chromedp.Flag("blink-settings", "imagesEnabled=true"),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var buf []byte
	screenshotPath := filepath.Join(screenshotDir, filename)

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.Sleep(3*time.Second),
	)

	if err == nil {
		err = chromedp.Run(ctx,
			chromedp.FullScreenshot(&buf, 90),
		)
	}

	if err != nil {
		return fmt.Errorf("screenshot alınamadı: %v", err)
	}

	if err := os.WriteFile(screenshotPath, buf, 0644); err != nil {
		return fmt.Errorf("screenshot kaydedilemedi: %v", err)
	}

	LogSuccess("Screenshot kaydedildi: %s (%d KB)", screenshotPath, len(buf)/1024)
	return nil
}

func TakeScreenshotSimple(url, name string, torPort int) bool {
	filename := sanitizeFilename(name) + "_" + time.Now().Format("20060102_150405") + ".png"

	LogInfo("Screenshot alınıyor: %s", name)

	err := TakeScreenshot(url, filename, torPort)
	if err != nil {
		LogWarning("Screenshot atlandı (%s): Site yavaş veya erişilemez", name)
		return false
	}

	return true
}
