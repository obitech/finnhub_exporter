package query

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/prometheus/client_golang/prometheus"
)

const promNamespace = "finnhub"

// Querier can send queries to the Finnhub.io API
type Querier interface {
	Do(context.Context, *finnhub.DefaultApiService, *prometheus.Registry,
		string) error
}
