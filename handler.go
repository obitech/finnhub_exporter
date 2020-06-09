package main

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"obitech/finnhub_exporter/endpoint"
)

// getStockID retrieves a StockID from URL values.
func getStockID(r *http.Request) (*endpoint.StockID, error) {
	symbol := r.URL.Query().Get(symbolParam)
	isin := r.URL.Query().Get(isinParam)
	cusip := r.URL.Query().Get(cusipParam)
	if symbol == "" && isin == "" && cusip == "" {
		return nil, fmt.Errorf("need to set parameter %q or %q or %q", symbolParam, isinParam, cusipParam)
	}
	return &endpoint.StockID{
		Symbol: symbol,
		ISIN: isin,
		CUSIP: cusip,
	}, nil
}

func queryHandler(w http.ResponseWriter, r *http.Request, log zap.Logger) {
	endpointName := r.URL.Query().Get(endpointParam)
	if endpointName == "" {
		http.Error(w, fmt.Sprintf("Need to pass %q parameter", endpointParam), http.StatusBadRequest)
		return
	}
}