package query

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/antihax/optional"
	"github.com/prometheus/client_golang/prometheus"
)

// CompanyProfile2 provides basic company company information.
type CompanyProfile2 struct{}

// Do makes a request against the /stock/profile2 endpoint.
func (c CompanyProfile2) Do(ctx context.Context,
	client *finnhub.DefaultApiService, registry *prometheus.Registry,
	symbol string) error {
	opts := &finnhub.CompanyProfile2Opts{
		Symbol: optional.NewString(symbol),
	}
	profile, _, err := client.CompanyProfile2(ctx, opts)
	if err != nil {
		return err
	}

	if profile.Name == "" {
		return fmt.Errorf("no data returned for stock: %s", symbol)
	}

	cp2Gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: promNamespace,
			Name:      "company_profile_2",
			Help:      "Displays general information of a company (free version of CompanyProfile)",
		},
		[]string{"symbol", "country", "currency", "exchange", "ipo",
			"marketCapitalization", "name", "shareOutstanding", "ticker",
			"weburl", "logo", "finnhubIndustry",
		},
	)
	registry.MustRegister(cp2Gauge)
	cp2Gauge.WithLabelValues(
		symbol, profile.Country, profile.Currency, profile.Exchange,
		profile.Ipo, strconv.FormatFloat(float64(profile.MarketCapitalization), 'g', -1, 32),
		profile.Name,
		strconv.FormatFloat(float64(profile.ShareOutstanding), 'g', -1, 32),
		profile.Ticker, profile.Weburl, profile.Logo, profile.FinnhubIndustry,
	).Set(1)

	return nil
}
