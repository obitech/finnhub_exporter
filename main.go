package main

import (
	"fmt"
	"os"
)

const (
	apiKeyEnv = "FINNHUB_API_KEY"
)

func main() {
	l, err := NewLogger("info")
	defer l.Sync()
	log := l.Sugar()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		log.Errorw("environment variable needs to be set",
			"env", apiKeyEnv,
		)
		os.Exit(1)
	}
}
