package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var logFile *os.File

// InitLogger initializes the log file
func InitLogger() {
	var err error
	logFile, err = os.OpenFile("go_discovery.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.Println("=== Go Discovery Plugin Started ===")
}

// CloseLogger ensures the log file is properly closed
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// OutputResponse logs and prints JSON output
func OutputResponse(data interface{}) {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Failed to encode output: %v", err)
	}
	log.Println("Response:", string(output))
	fmt.Println(string(output))
}
