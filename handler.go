package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	"obitech/finnhub_exporter/endpoint"
)

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

		// TODO: Adjust with Prometheus Scrape-Timout Header
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		_, auth := NewFinnhubClient(ctx, apiKey)
		r = r.WithContext(auth)

		_, err := getStockID(r)
		if err != nil {
			log.Debug("Unable to extract stock id from URL values", zap.Error(err))
			http.Error(w, fmt.Sprintf("Unable to extract stock id: %v", err), http.StatusBadRequest)
			return
		}

		// We don't want to unit test the actual API call against Finnhub.io
		if test {
			return
		}
	}

}
