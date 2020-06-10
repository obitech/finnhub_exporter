package query

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/prometheus/client_golang/prometheus"
)

const promNamespace = "finnhub"

// Querier can send queries to the Finnhub.io API
type Querier interface {
	// Do performs a query against a specific endpoint implemented by the
	// caller type.
	Do(context.Context, *finnhub.DefaultApiService, *prometheus.Registry,
		string) error
}
