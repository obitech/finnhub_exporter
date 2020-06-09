package endpoint

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/prometheus/client_golang/prometheus"
)

type RequestFn interface {
	Do(finnhub.DefaultApiService, context.Context, *prometheus.Registry) error
}

type CompanyProfile2 finnhub.CompanyProfile2

func (c CompanyProfile2) Do(client finnhub.DefaultApiService, ctx context.Context, id StockID, registry *prometheus.Registry) error {
	return nil
}

// StockID identifies a stock by either a Symbol, ISIN or CUSIP.
type StockID struct {
	Symbol string
	ISIN   string
	CUSIP  string
}

