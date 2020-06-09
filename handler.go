package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"obitech/finnhub_exporter/endpoint"
)

const promNamespace = "finnhub"

var modules = map[string]endpoint.RequestFn{
	"companyprofile2": endpoint.CompanyProfile2{},
}

// getStockID retrieves a StockID from URL values.
func getStockID(r *http.Request) (*endpoint.StockID, error) {
	symbol := r.URL.Query().Get(symbolParam)
	isin := r.URL.Query().Get(isinParam)
	cusip := r.URL.Query().Get(cusipParam)
	if symbol == "" && isin == "" && cusip == "" {
		return nil, fmt.Errorf("need to set URL parameter %q or %q or %q", symbolParam, isinParam, cusipParam)
	}
	return &endpoint.StockID{
		Symbol: symbol,
		ISIN:   isin,
		CUSIP:  cusip,
	}, nil
}

func queryHandler(apiKey string, log *zap.Logger, test bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpointName := r.URL.Query().Get(endpointParam)
		if endpointName == "" {
			http.Error(w, fmt.Sprintf("Missing URL parameter %q", endpointParam), http.StatusBadRequest)
			return
		}

		ep, ok := modules[endpointName]
		if !ok {
			http.Error(w, fmt.Sprintf("Endpoint %q is not supported", endpointParam), http.StatusBadRequest)
			return
		}

		// TODO: Adjust with Prometheus Scrape-Timout Header
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		client, auth := NewFinnhubClient(ctx, apiKey)
		r = r.WithContext(auth)

		stockID, err := getStockID(r)
		if err != nil {
			log.Debug("Unable to extract stock id from URL values", zap.Error(err))
			http.Error(w, fmt.Sprintf("Unable to extract stock id: %v", err), http.StatusBadRequest)
			return
		}

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
		err = ep.Do(auth, client, registry, stockID)
		duration := time.Since(start).Seconds()
		queryDurationGauge.Set(duration)
		if err != nil {
			querySuccessGauge.Set(0)
			log.Info("query to finnhub failed",
				zap.Error(err),
				zap.Duration("query_duration", time.Duration(duration)),
			)
		} else {
			querySuccessGauge.Set(1)
			log.Info("query to finnhub successful",
				zap.String("endpoint", endpointName),
				zap.String("stockID", stockID.String()),
				zap.Duration("query_duration", time.Duration(duration)),
			)
		}

		ph := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		ph.ServeHTTP(w, r)
	}
}
