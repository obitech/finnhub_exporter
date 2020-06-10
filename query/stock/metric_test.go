package stock

import (
	"testing"

	"github.com/magiconair/properties/assert"
)

func TestCleanName(t *testing.T) {
	var data = []struct {
		name     string
		in       string
		expected string
	}{
		{
			name:     "empty string",
			in:       "",
			expected: "",
		},
		{
			name:     "normal string",
			in:       "test",
			expected: "test",
		},
		{
			name:     "only slash",
			in:       "/",
			expected: "_divided_",
		},
		{
			name:     "only ampersand",
			in:       "&",
			expected: "_and_",
		},
		{
			name:     "metric name with dash",
			in:       "currentEv/freeCashFlowAnnual",
			expected: "currentEv_divided_freeCashFlowAnnual",
		},
		{
			name:     "metric name with ampersand",
			in:       "priceRelativeToS&P50013Week",
			expected: "priceRelativeToS_and_P50013Week",
		},
		{
			name:     "metric name with dash and ampersand",
			in:       "currentEv/freeCashFlowAnnual_priceRelativeToS&P50013Week",
			expected: "currentEv_divided_freeCashFlowAnnual_priceRelativeToS_and_P50013Week",
		},
	}

	for _, test := range data {
		t.Run(test.name, func(t *testing.T) {
			actual := cleanName(test.in)
			assert.Equal(t, actual, test.expected)
		})
	}
}
