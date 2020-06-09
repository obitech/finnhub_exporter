package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"obitech/finnhub_exporter/endpoint"
)

func TestGetStockID(t *testing.T) {
	data := []struct {
		name     string
		params   url.Values
		expected *endpoint.StockID
		err      bool
	}{
		{
			name: "empty params",
			err:  true,
		},
		{
			name:   "only symbol",
			params: url.Values{"symbol": []string{"AAPL"}},
			expected: &endpoint.StockID{
				Symbol: "AAPL",
			},
		},
		{
			name:   "only cusip",
			params: url.Values{"cusip": []string{"abc"}},
			expected: &endpoint.StockID{
				CUSIP: "abc",
			},
		},
		{
			name:   "only isin",
			params: url.Values{"isin": []string{"123"}},
			expected: &endpoint.StockID{
				ISIN: "123",
			},
		},
		{
			name: "symbol,isin,cusip",
			params: url.Values{
				"isin":   []string{"123"},
				"symbol": []string{"AAPL"},
				"cusip":  []string{"abc"},
			},
			expected: &endpoint.StockID{
				ISIN:   "123",
				Symbol: "AAPL",
				CUSIP:  "abc",
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

func TestQueryHandler(t *testing.T) {
	log := zap.NewNop()
	var data = []struct {
		name           string
		path           string
		apiKey         string
		log            *zap.Logger
		expectedStatus int
	}{
		{
			name:           "missing endpoint param",
			path:           "/test",
			log:            log,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "endpoint not supported",
			path:           "/test?endpoint=abc",
			log:            log,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "endpoint present, symbol passed",
			path:           "/test?endpoint=companyprofile2&symbol=AAPL",
			log:            log,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "endpoint present, ISIN passed",
			path:           "/test?endpoint=companyprofile2&isin=123",
			log:            log,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "endpoint present, cusip passed",
			path:           "/test?endpoint=companyprofile2&cusip=abc",
			log:            log,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "endpoint present, all stock ids passed",
			path:           "/test?endpoint=companyprofile2&cusip=abc&symbol=AAPL&isin=123",
			log:            log,
			expectedStatus: http.StatusOK,
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", test.path, nil)
			assert.NoError(t, err, "creating a test request shouldn't fail")

			rr := httptest.NewRecorder()
			handler := QueryHandler(test.apiKey, test.log, true)
			handler.ServeHTTP(rr, req)

			if assert.Equal(t, test.expectedStatus, rr.Code, rr.Body) {
				t.Logf("Body: %q", rr.Body.String())
			}
		})
	}
}
