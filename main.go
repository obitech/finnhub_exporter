package main

import (
	"fmt"
	"os"
)

const (
	apiKeyEnv = "FINNHUB_API_KEY"
)

func main() {
	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		fmt.Printf("environment variable %q needs to be set\n", apiKeyEnv)
		os.Exit(1)
	}
}
