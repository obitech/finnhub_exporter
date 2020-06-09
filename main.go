package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"github.com/spf13/cobra"
)

const (
	apiKeyEnv     = "FINNHUB_API_KEY"
	endpointParam = "endpoint"
	symbolParam   = "symbol"
	isinParam     = "isin"
	cusipParam    = "cusip"
)

var (
	address = ":9780"
	rootCmd = &cobra.Command{
		Use:   "run",
		Short: "Export financial data from finnhub.io",
		Run:   run,
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

	// TODO: Adjust with Prometheus Scrape-Timout Header
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(30*time.Second))
	defer cancel()
	client, auth := NewFinnhubClient(ctx, apiKey)
	opts := &finnhub.CompanyProfile2Opts{
		Symbol: optional.NewString("AAPL"),
	}
	profile, _, err := client.CompanyProfile2(auth, opts)
	if err != nil {
		log.Errorw("unable to get CompanyProfile2",
			"error", err,
		)
	}
	fmt.Println(profile)
}

func init() {
	f := rootCmd.Flags()
	f.StringVarP(&address, "web.listen-address", "a", address, "The address to listen on for HTTP requests.")
}

func main() {
	Execute()
}
