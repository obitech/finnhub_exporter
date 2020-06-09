package endpoint

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/prometheus/client_golang/prometheus"
)

type RequestFn interface {
	Do(context.Context, finnhub.DefaultApiService, *prometheus.Registry, StockID) error
}

type CompanyProfile2 finnhub.CompanyProfile2

func (c CompanyProfile2) Do(ctx context.Context, client finnhub.DefaultApiService, registry *prometheus.Registry, id StockID) error {
	return nil
}

// StockID identifies a stock by either a Symbol, ISIN or CUSIP.
type StockID struct {
	Symbol string
	ISIN   string
	CUSIP  string
}
