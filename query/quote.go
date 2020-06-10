package query

import (
	"context"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/prometheus/client_golang/prometheus"
)

// Quote wraps finnhub.Quote. Gets quote data for stocks.
type Quote struct {
	finnhub.Quote
}

// Do makes a request against the /quote endpoint.
func (c Quote) Do(ctx context.Context,
	client *finnhub.DefaultApiService, registry *prometheus.Registry,
	symbol string) error {
	subsystem := "quote"
	gaugeLabels := []string{"symbol"}

	openGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: subsystem,
		Name:      "open",
	}, gaugeLabels)
	highGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: subsystem,
		Name:      "high",
	}, gaugeLabels)
	lowGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: subsystem,
		Name:      "low",
	}, gaugeLabels)
	currentGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: subsystem,
		Name:      "current",
	}, gaugeLabels)
	prevCloseGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: subsystem,
		Name:      "prev_close",
	}, gaugeLabels)
	registry.MustRegister(openGauge)
	registry.MustRegister(highGauge)
	registry.MustRegister(lowGauge)
	registry.MustRegister(currentGauge)
	registry.MustRegister(prevCloseGauge)

	quote, _, err := client.Quote(ctx, symbol)
	if err != nil {
		return err
	}

	openGauge.WithLabelValues(symbol).Set(float64(quote.O))
	highGauge.WithLabelValues(symbol).Set(float64(quote.H))
	lowGauge.WithLabelValues(symbol).Set(float64(quote.L))
	currentGauge.WithLabelValues(symbol).Set(float64(quote.C))
	prevCloseGauge.WithLabelValues(symbol).Set(float64(quote.Pc))

	return nil
}
