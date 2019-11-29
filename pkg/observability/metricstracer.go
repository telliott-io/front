package observability

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var _ opentracing.Tracer = &metricsTracer{}

func NewMetricsTracer(namespace string, tracer opentracing.Tracer) opentracing.Tracer {
	return &metricsTracer{
		Tracer:     tracer,
		namespace:  namespace,
		counters:   make(map[string]prometheus.Counter),
		histograms: make(map[string]prometheus.Histogram),
	}
}

type metricsTracer struct {
	opentracing.Tracer

	namespace string

	counterMutex sync.Mutex
	counters     map[string]prometheus.Counter

	histogramMutex sync.Mutex
	histograms     map[string]prometheus.Histogram
}

func normalizeMetricName(name string) string {
	return strings.ReplaceAll(
		strings.ReplaceAll(name, "-", "_"),
		"/", "_",
	)
}

func (m *metricsTracer) getHistogram(name string) prometheus.Histogram {
	m.histogramMutex.Lock()
	defer m.histogramMutex.Unlock()
	if histogram, exists := m.histograms[name]; exists {
		return histogram
	}
	histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    fmt.Sprintf("%v_duration_ms", normalizeMetricName(name)),
		Buckets: prometheus.ExponentialBuckets(10, 10, 5),
	})
	prometheus.Register(histogram)

	m.histograms[name] = histogram
	return histogram
}

func (m *metricsTracer) getCounter(name string) prometheus.Counter {
	m.counterMutex.Lock()
	defer m.counterMutex.Unlock()
	if counter, exists := m.counters[name]; exists {
		return counter
	}
	counter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: m.namespace,
		Name:      fmt.Sprintf("%v_count", normalizeMetricName(name)),
	})
	m.counters[name] = counter
	return counter
}

func (m *metricsTracer) StartSpan(operationName string, opts ...opentracing.StartSpanOption) opentracing.Span {
	spanCount := m.getCounter(operationName)
	spanCount.Inc()

	span := m.Tracer.StartSpan(operationName, opts...)
	return newSpan(operationName, span, m)
}

func newSpan(operationName string, inner opentracing.Span, tracer *metricsTracer) *spanWrapper {
	return &spanWrapper{
		name:   operationName,
		Span:   inner,
		start:  time.Now(),
		tracer: tracer,
	}
}

type spanWrapper struct {
	opentracing.Span

	start time.Time

	name string

	tracer *metricsTracer
}

func (s *spanWrapper) Finish() {
	s.Span.Finish()

	hg := s.tracer.getHistogram(s.name)
	hg.Observe(float64(time.Since(s.start) / time.Millisecond))
}

func (s *spanWrapper) FinishWithOptions(opts opentracing.FinishOptions) {
	s.Span.FinishWithOptions(opts)
}
