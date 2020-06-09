package main

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go"
)

func NewFinnhubClient(ctx context.Context, apiKey string) (*finnhub.DefaultApiService, context.Context) {
	client := finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi
	auth := context.WithValue(ctx, finnhub.ContextAPIKey, finnhub.APIKey{
		Key: apiKey,
	})

	return client, auth
}
