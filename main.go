package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	apiKeyEnv = "FINNHUB_API_KEY"
)

var (
	address  = ":9780"
	logLevel = "info"
	rootCmd  = &cobra.Command{
		Use:   "run",
		Short: "Export financial data from finnhub.io",
		RunE:  run,
	}
)

func init() {
	f := rootCmd.Flags()
	f.StringVarP(&address, "web.listen-address", "a", address, "The address to listen on for HTTP requests.")
	f.StringVarP(&logLevel, "log.level", "l", logLevel, "log level (debug, info, warn, error). Empty or invalid values will fallback to info")
}

func main() {
	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	log, err := newLogger(logLevel)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	apiKey := os.Getenv(apiKeyEnv)
	if apiKey == "" {
		log.Error("environment variable missing",
			zap.String("env", apiKeyEnv),
		)
		os.Exit(1)
	}

	http.HandleFunc("/query", QueryHandler(apiKey, log, false))
	http.Handle("/metrics", promhttp.Handler())

	srv := http.Server{Addr: address}
	srvc := make(chan struct{})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("Server listening", zap.String("address", address))
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Error("Error starting HTTP server", zap.Error(err))
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			log.Info("Received SIGTERM, exiting gracefully...")
			return nil
		case <-srvc:
			log.Error("Aborting")
			return fmt.Errorf("aborting")
		}
	}
}
