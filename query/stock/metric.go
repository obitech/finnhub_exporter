package stock

import (
	"context"
	"fmt"
	"strings"

	"github.com/Finnhub-Stock-API/finnhub-go"
	"github.com/prometheus/client_golang/prometheus"

	"obitech/finnhub_exporter/config"
)

// Metric gets company basic financials such as margin, P/E ratio,
// 52-week high/low etc. It only queries
type Metric struct{}

// Do makes a request against the /stock/metric endpoint with the metric=all
// parameter.
func (q Metric) Do(ctx context.Context,
	client *finnhub.DefaultApiService, registry *prometheus.Registry,
	symbol string) error {

	metrics, _, err := client.CompanyBasicFinancials(ctx, symbol, "all")
	if err != nil {
		return err
	}

	if len(metrics.Metric) == 0 {
		return fmt.Errorf("invalid symbol: %s", symbol)
	}

	for k, v := range metrics.Metric {
		switch t := v.(type) {
		case float32:
			g := newMetricGauge(cleanName(k), symbol)
			registry.MustRegister(g)
			g.Set(float64(t))
		case float64:
			g := newMetricGauge(cleanName(k), symbol)
			registry.MustRegister(g)
			g.Set(t)
		case int32:
			g := newMetricGauge(cleanName(k), symbol)
			registry.MustRegister(g)
			g.Set(float64(t))
		case int64:
			g := newMetricGauge(cleanName(k), symbol)
			registry.MustRegister(g)
			g.Set(float64(t))
		default:
			continue
		}
	}
	return nil
}

// newMetricGauge creates a new Prometheus Gauge for a stock metric.
func newMetricGauge(name, symbol string) prometheus.Gauge {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: config.PromNamespace,
		Subsystem: "stock_metric",
		Name:      name,
	}, []string{config.SymbolLabel}).WithLabelValues(symbol)
}

// cleanName removes invalid characters in a Prometheus metric name.
func cleanName(name string) string {
	if !strings.Contains(name, "/") && !strings.Contains(name, "&") {
		return name
	}
	name = strings.ReplaceAll(name, "/", "_divided_")
	return strings.ReplaceAll(name, "&", "_and_")
}
