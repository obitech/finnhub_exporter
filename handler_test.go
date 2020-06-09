package main

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"obitech/finnhub_exporter/endpoint"
)

func TestGetStockID(t *testing.T) {
	data := []struct {
		name     string
		params  url.Values
		expected *endpoint.StockID
		err      bool
	}{
		{
			name:     "empty params",
			err:      true,
		},
		{
			name: "only symbol",
			params: url.Values{"symbol": []string{"AAPL"}},
			expected: &endpoint.StockID{
				Symbol: "AAPL",
			},
		},
		{
			name: "only cusip",
			params: url.Values{"cusip": []string{"abc"}},
			expected: &endpoint.StockID{
				CUSIP: "abc",
			},
		},
		{
			name: "only isin",
			params: url.Values{"isin": []string{"123"}},
			expected: &endpoint.StockID{
				ISIN: "123",
			},
		},
		{
			name: "symbol,isin,cusip",
			params: url.Values{
				"isin": []string{"123"},
				"symbol": []string{"AAPL"},
				"cusip": []string{"abc"},
			},
			expected: &endpoint.StockID{
				ISIN: "123",
				Symbol: "AAPL",
				CUSIP: "abc",
			},
		},

	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			var req *http.Request
			if test.params != nil {
				req, _ = http.NewRequest("GET", "/test?"+test.params.Encode(), nil)
			} else {
				req, _ = http.NewRequest("GET", "/test", nil)
			}

			actual, err := getStockID(req)
			if err != nil {
				if !test.err {
					t.Fatalf("unexpected error: %v", err)
				}
			}
			assert.Equal(t, actual, test.expected, "should be equal")
		})
	}
}
