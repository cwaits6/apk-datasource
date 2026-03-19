package metrics

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

var (
	attrMethod     = attribute.Key("method")
	attrPath       = attribute.Key("path")
	attrStatusCode = attribute.Key("status_code")
	attrStatus     = attribute.Key("status")
)

// Metrics holds all OTel instruments for the application.
type Metrics struct {
	httpRequestsTotal   metric.Int64Counter
	httpRequestDuration metric.Float64Histogram
	refreshTotal        metric.Int64Counter
	refreshDuration     metric.Float64Histogram
	refreshPackages     metric.Int64Gauge
	serverReady         metric.Int64Gauge
}

// Setup creates a Prometheus exporter, registers an OTel MeterProvider, and
// returns the Metrics instruments plus an http.Handler for /metrics.
func Setup() (*Metrics, http.Handler, error) {
	exporter, err := prometheus.New()
	if err != nil {
		return nil, nil, err
	}

	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	otel.SetMeterProvider(provider)

	meter := provider.Meter("apk-datasource")
	m := &Metrics{}

	m.httpRequestsTotal, err = meter.Int64Counter("http_requests_total",
		metric.WithDescription("Total HTTP requests"))
	if err != nil {
		return nil, nil, err
	}

	m.httpRequestDuration, err = meter.Float64Histogram("http_request_duration_seconds",
		metric.WithDescription("HTTP request duration in seconds"))
	if err != nil {
		return nil, nil, err
	}

	m.refreshTotal, err = meter.Int64Counter("refresh_total",
		metric.WithDescription("Total index refreshes"))
	if err != nil {
		return nil, nil, err
	}

	m.refreshDuration, err = meter.Float64Histogram("refresh_duration_seconds",
		metric.WithDescription("Index refresh duration in seconds"))
	if err != nil {
		return nil, nil, err
	}

	m.refreshPackages, err = meter.Int64Gauge("refresh_packages",
		metric.WithDescription("Number of packages after last refresh"))
	if err != nil {
		return nil, nil, err
	}

	m.serverReady, err = meter.Int64Gauge("server_ready",
		metric.WithDescription("Whether the server is ready (0 or 1)"))
	if err != nil {
		return nil, nil, err
	}

	return m, promhttp.Handler(), nil
}

// Noop returns a Metrics with nil instruments. All helper methods are nil-safe.
func Noop() *Metrics {
	return &Metrics{}
}

// RecordHTTPRequest records an HTTP request's method, path, status code, and duration.
func (m *Metrics) RecordHTTPRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
	if m == nil || m.httpRequestsTotal == nil {
		return
	}
	attrs := metric.WithAttributes(
		attrMethod.String(method),
		attrPath.String(path),
		attrStatusCode.String(strconv.Itoa(statusCode)),
	)
	m.httpRequestsTotal.Add(ctx, 1, attrs)
	m.httpRequestDuration.Record(ctx, duration.Seconds(), attrs)
}

// RecordRefresh records a refresh attempt's status and duration.
func (m *Metrics) RecordRefresh(ctx context.Context, status string, duration time.Duration, packageCount int64) {
	if m == nil || m.refreshTotal == nil {
		return
	}
	attrs := metric.WithAttributes(attrStatus.String(status))
	m.refreshTotal.Add(ctx, 1, attrs)
	m.refreshDuration.Record(ctx, duration.Seconds(), attrs)
	m.refreshPackages.Record(ctx, packageCount)
}

// SetReady sets the server readiness gauge.
func (m *Metrics) SetReady(ctx context.Context, ready bool) {
	if m == nil || m.serverReady == nil {
		return
	}
	var v int64
	if ready {
		v = 1
	}
	m.serverReady.Record(ctx, v)
}
