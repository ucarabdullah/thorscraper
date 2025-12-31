package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var logFile *os.File

func InitLogger() error {
	outputDir := "logs"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	logPath := filepath.Join(outputDir, "scan_report.log")

	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// Başlık
	writeLog("========================================")
	writeLog("TOR SCRAPER - DETAYLI LOG KAYDI")
	writeLog("========================================")
	writeLog("Başlangıç: %s", time.Now().Format("2006-01-02 15:04:05"))
	writeLog("Log Dosyası: %s", logPath)
	writeLog("========================================\n")

	return nil
}

func CloseLogger() {
	if logFile != nil {
		writeLog("\n========================================")
		writeLog("Bitiş: %s", time.Now().Format("2006-01-02 15:04:05"))
		writeLog("========================================")
		logFile.Close()
	}
}

func writeLog(format string, args ...interface{}) {
	if logFile != nil {
		msg := fmt.Sprintf(format, args...)
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(logFile, "[%s] %s\n", timestamp, msg)
	}
}

func LogInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	writeLog("[INFO] %s", msg)
}

func LogSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	writeLog("[SUCCESS] %s", msg)
}

func LogError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	writeLog("[ERROR] %s", msg)
}

func LogWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	writeLog("[WARNING] %s", msg)
}

func LogDebug(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	writeLog("[DEBUG] %s", msg)
}

func LogStep(step int, total int, message string) {
	writeLog("[STEP %d/%d] %s", step, total, message)
}

func PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("ℹ️  %s\n", msg)
	LogInfo(msg)
}

func PrintSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("✓ %s\n", msg)
	LogSuccess(msg)
}

func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("✗ %s\n", msg)
	LogError(msg)
}

func PrintWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf("⚠  %s\n", msg)
	LogWarning(msg)
}

func PrintStep(step int, total int, message string) {
	fmt.Printf("\n[%d/%d] %s\n", step, total, message)
	LogStep(step, total, message)
}

func PrintDivider() {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}
