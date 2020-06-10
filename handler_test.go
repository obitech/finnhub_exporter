package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

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
			name:           "endpoint present, symbol missing",
			path:           "/test?endpoint=companyprofile2",
			log:            log,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "endpoint present, symbol passed",
			path:           "/test?endpoint=companyprofile2&symbol=AAPL",
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
