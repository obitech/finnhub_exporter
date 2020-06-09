package endpoint

import (
	"context"
	"fmt"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"github.com/prometheus/client_golang/prometheus"
)

type RequestFn interface {
	Do(context.Context, *finnhub.DefaultApiService, *prometheus.Registry, *StockID) error
}

type CompanyProfile2 finnhub.CompanyProfile2

func (c CompanyProfile2) Do(ctx context.Context, client *finnhub.DefaultApiService, registry *prometheus.Registry, id *StockID) error {
	opts := &finnhub.CompanyProfile2Opts{
		Symbol: optional.NewString(id.Symbol),
		Isin:   optional.NewString(id.ISIN),
		Cusip:  optional.NewString(id.CUSIP),
	}
	_, _, err := client.CompanyProfile2(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

// StockID identifies a stock by either a Symbol, ISIN or CUSIP.
type StockID struct {
	Symbol string
	ISIN   string
	CUSIP  string
}

func (id StockID) String() string {
	return fmt.Sprintf("Symbol: %s, ISIN: %s, CUSIP: %s", id.Symbol, id.ISIN, id.CUSIP)
}
