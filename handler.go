package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"obitech/finnhub_exporter/query"
	"obitech/finnhub_exporter/query/stock"
)

const (
	promNamespace = "finnhub"
	endpointParam = "endpoint"
	symbolParam   = "symbol"
)

var modules = map[string]query.Querier{
	"companyprofile2": stock.CompanyProfile2{},
	"quote":           query.Quote{},
	"metric":          stock.Metric{},
}

// QueryHandler defines an uninstrumented Prometheus handler that allows for
// dynamic queries against the Finnhub.io API.
// Pass the wished endpoint and stock symbol to make a query.
func QueryHandler(apiKey string, log *zap.Logger, test bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointName := r.URL.Query().Get(endpointParam)
		if endpointName == "" {
			http.Error(w, fmt.Sprintf("URL parameter missing: %s",
				endpointParam), http.StatusBadRequest)
			return
		}

		ep, ok := modules[endpointName]
		if !ok {
			http.Error(w, fmt.Sprintf("Endpoint not supported: %s",
				endpointParam), http.StatusBadRequest)
			return
		}

		symbol := r.URL.Query().Get(symbolParam)
		if symbol == "" {
			http.Error(w, fmt.Sprintf("URL parameter missing: %s",
				symbolParam), http.StatusBadRequest)
			return
		}

		// TODO: Adjust with Prometheus Scrape-Timout Header
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		client, auth := query.NewFinnhubClient(ctx, apiKey)
		r = r.WithContext(auth)

		// We don't want to unit test the actual API call against Finnhub.io
		if test {
			return
		}
		registry := prometheus.NewRegistry()
		querySuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: promNamespace,
			Name:      "query_success",
			Help:      "Displays whether a query to the Finnhub API was successful",
		})
		queryDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: promNamespace,
			Name:      "query_duration",
			Help:      "Returns how long a query to the Finnhub API took to complete in seconds",
		})
		registry.MustRegister(queryDurationGauge)
		registry.MustRegister(querySuccessGauge)

		start := time.Now()
		err := ep.Do(auth, client, registry, symbol)
		duration := time.Since(start)
		queryDurationGauge.Set(duration.Seconds())
		if err != nil {
			querySuccessGauge.Set(0)
			log.Error("query to finnhub failed",
				zap.Error(err),
				zap.Duration("query_duration", duration),
			)
		} else {
			querySuccessGauge.Set(1)
			log.Info("query to finnhub successful",
				zap.String("endpoint", endpointName),
				zap.String("symbol", symbol),
				zap.Duration("query_duration", duration),
			)
		}

		ph := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		ph.ServeHTTP(w, r)
	}
}
