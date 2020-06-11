package query

import (
	"context"
	"fmt"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/prometheus/client_golang/prometheus"

	"obitech/finnhub_exporter/config"
)

// Quote gets quote data for stocks.
type Quote struct{}

// Do makes a request against the /quote endpoint.
func (c Quote) Do(ctx context.Context,
	client *finnhub.DefaultApiService, registry *prometheus.Registry,
	symbol string) error {
	quote, _, err := client.Quote(ctx, symbol)
	if err != nil {
		return err
	}

	// Check if symbol exists.
	if quote.O == float32(0) && quote.C == float32(0) && quote.Pc == 0 {
		return fmt.Errorf("symbol doesn't exists: %s", symbol)
	}

	subsystem := "quote"
	gaugeLabels := []string{config.SymbolLabel}

	openGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: promNamespace,
		Subsystem: subsystem,
		Name:      "open",
	}, gaugeLabels)
	highGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: config.PromNamespace,
		Subsystem: subsystem,
		Name:      "high",
	}, gaugeLabels)
	lowGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: config.PromNamespace,
		Subsystem: subsystem,
		Name:      "low",
	}, gaugeLabels)
	currentGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: config.PromNamespace,
		Subsystem: subsystem,
		Name:      "current",
	}, gaugeLabels)
	prevCloseGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: config.PromNamespace,
		Subsystem: subsystem,
		Name:      "prev_close",
	}, gaugeLabels)
	registry.MustRegister(openGauge)
	registry.MustRegister(highGauge)
	registry.MustRegister(lowGauge)
	registry.MustRegister(currentGauge)
	registry.MustRegister(prevCloseGauge)

	openGauge.WithLabelValues(symbol).Set(float64(quote.O))
	highGauge.WithLabelValues(symbol).Set(float64(quote.H))
	lowGauge.WithLabelValues(symbol).Set(float64(quote.L))
	currentGauge.WithLabelValues(symbol).Set(float64(quote.C))
	prevCloseGauge.WithLabelValues(symbol).Set(float64(quote.Pc))

	return nil
}
