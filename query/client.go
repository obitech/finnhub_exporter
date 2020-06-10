package query

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go"
)

// NewFinnhubClient returns a finnhub.
// DefaultApiService and a new context with authentication parameters present.
func NewFinnhubClient(ctx context.Context, apiKey string) (
	*finnhub.DefaultApiService, context.Context) {
	client := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(ctx, finnhub.ContextAPIKey, finnhub.APIKey{
		Key: apiKey,
	})

	return client, auth
}
