package metricshdl

import (
	"context"
	"fmt"
	"github.com/adrianoccosta/exercise-qonto/internal/middleware"
	"github.com/adrianoccosta/exercise-qonto/tools"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

const (
	// KeyPath is the key for request path
	KeyPath MetricsKey = "handler"
	// KeyMethod is the key for request method
	KeyMethod MetricsKey = "method"
	// KeyPartner is the key for the requesting partner
	KeyPartner MetricsKey = "partner"
	// KeyStatus is the key for response statusCode code
	KeyStatus MetricsKey = "status"
)

var (
	defaultMetricKeys = []MetricsKey{
		KeyPath,
		KeyMethod,
		KeyStatus,
		KeyPartner,
	}
)

// MetricsConfig .
type MetricsConfig struct {
	Keys                   []MetricsKey
	Namespace              string
	Subsystem              string
	CounterName            string
	CounterHelp            string
	HistogramName          string
	HistogramHelp          string
	GaugeName              string
	GaugeHelp              string
	ContextMetricsBeforeFn func(r *http.Request) context.Context
	ContextMetricsAfterFn  func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context
}

// MetricsKey defines a type to configure the possible values for the metrics middleware
type MetricsKey string

// Metrics defines the interface of the service for middleware usage, including registry into a prometheus registry
type Metrics interface {
	Handler(http.Handler) http.Handler
	Registry() prometheus.Gatherer
}

// DefaultMetricKeys returns the default metric keys
func DefaultMetricKeys() []MetricsKey {
	return defaultMetricKeys
}

// DefaultContextMetricsBeforeFunction returns the default context metrics before function
func DefaultContextMetricsBeforeFunction() func(r *http.Request) context.Context {
	return defaultContextMetricsBefore
}

// DefaultContextMetricsAfterFunction returns the default context metrics after function
func DefaultContextMetricsAfterFunction() func(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {
	return defaultContextMetricsAfter
}

// NewMetrics .
func NewMetrics(namespace, subsystem string) Metrics {
	return NewMetricsWithConfig(MetricsConfig{
		Keys:                   defaultMetricKeys,
		Namespace:              namespace,
		Subsystem:              subsystem,
		CounterName:            "http_requests_count",
		CounterHelp:            "Number of requests received",
		HistogramName:          "http_requests_duration",
		HistogramHelp:          "Request Duration in microseconds",
		ContextMetricsBeforeFn: defaultContextMetricsBefore,
		ContextMetricsAfterFn:  defaultContextMetricsAfter,
	})
}

// NewMetricsWithConfig .
func NewMetricsWithConfig(cfg MetricsConfig) Metrics {
	labels := toSliceOfStrings(cfg.Keys)
	m := metrics{
		cfg: cfg,
	}
	var prometheusCollectors []prometheus.Collector
	if cfg.CounterName != "" {
		m.counterMetrics = prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.CounterName,
			Help:      cfg.CounterHelp,
		}, labels)
		prometheusCollectors = append(prometheusCollectors, m.counterMetrics)
	}

	if cfg.HistogramName != "" {
		m.histogramMetrics = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.HistogramName,
			Help:      cfg.HistogramHelp,
		}, labels)
		prometheusCollectors = append(prometheusCollectors, m.histogramMetrics)
	}

	if cfg.GaugeName != "" {
		m.gaugeMetrics = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: cfg.Namespace,
			Subsystem: cfg.Subsystem,
			Name:      cfg.GaugeName,
			Help:      cfg.GaugeHelp,
		}, labels)
		prometheusCollectors = append(prometheusCollectors, m.gaugeMetrics)
	}

	prometheusRegistry := prometheus.NewRegistry()
	prometheusRegistry.MustRegister(prometheusCollectors...)

	m.registry = prometheusRegistry

	return &m
}

type metrics struct {
	counterMetrics   *prometheus.CounterVec
	histogramMetrics *prometheus.HistogramVec
	gaugeMetrics     *prometheus.GaugeVec
	cfg              MetricsConfig
	registry         prometheus.Gatherer
}

func (m metrics) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := time.Now()

		ctx := r.Context()
		if m.cfg.ContextMetricsBeforeFn != nil {
			ctx = m.cfg.ContextMetricsBeforeFn(r)
		}

		lrw := middleware.WrapResponseWriter(w, nil)
		next.ServeHTTP(lrw, r.WithContext(ctx))

		ctx = context.WithValue(ctx, KeyStatus, lrw.Status())
		if m.cfg.ContextMetricsAfterFn != nil {
			ctx = m.cfg.ContextMetricsAfterFn(ctx, w, r)
		}

		m.registerMetrics(ctx, begin)

	})
}

func (m metrics) Registry() prometheus.Gatherer {
	return m.registry
}

func (m metrics) registerMetrics(ctx context.Context, begin time.Time) {
	labels := toPrometheusLabels(ctx, m.cfg.Keys)
	if m.counterMetrics != nil {
		m.counterMetrics.With(labels).Inc()
	}

	if m.gaugeMetrics != nil {
		m.gaugeMetrics.With(labels).Set(time.Since(begin).Seconds())
	}

	if m.histogramMetrics != nil {
		m.histogramMetrics.With(labels).Observe(time.Since(begin).Seconds())
	}
}

func partnerURN(r *http.Request) string {
	val := r.Header.Get(tools.HeaderXPartnerURN)
	if val != "" {
		return val
	}

	return "unknown"
}

func toPrometheusLabels(ctx context.Context, keys []MetricsKey) prometheus.Labels {
	res := prometheus.Labels{}
	for _, k := range keys {
		res[string(k)] = fmt.Sprint(ctx.Value(k))
	}
	return res
}

func defaultContextMetricsBefore(r *http.Request) context.Context {
	ctx := context.WithValue(r.Context(), KeyPath, r.URL.Path)

	route := mux.CurrentRoute(r)
	if route != nil {
		if pathTemplate, err := route.GetPathTemplate(); err == nil && pathTemplate != "" {
			ctx = context.WithValue(r.Context(), KeyPath, pathTemplate)
		}
	}

	ctx = context.WithValue(ctx, KeyMethod, r.Method)
	ctx = context.WithValue(ctx, KeyPartner, partnerURN(r))

	return ctx
}

func defaultContextMetricsAfter(ctx context.Context, _ http.ResponseWriter, _ *http.Request) context.Context {
	// do nothing by default
	return ctx
}

func toSliceOfStrings(keys []MetricsKey) (res []string) {
	for _, k := range keys {
		res = append(res, string(k))
	}
	return res
}
