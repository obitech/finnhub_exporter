package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	apiKeyEnv = "FINNHUB_API_KEY"
)

var (
	address = ":9780"
	rootCmd = &cobra.Command{
		Use: "run",
		Short: "Export financial data from finnhub.io",
		Run: run,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	l, err := NewLogger("info")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer l.Sync()
	log := l.Sugar()

	apiKey := os.Getenv(apiKeyEnv)

	if apiKey == "" {
		log.Errorw("environment variable needs to be set",
			"env", apiKeyEnv,
		)
		os.Exit(1)
	}
}

func init(){
	f := rootCmd.Flags()
	f.StringVarP(&address, "web.listen-address", "a", address, "The address to listen on for HTTP requests.")
}

func main() {
	Execute()
}
