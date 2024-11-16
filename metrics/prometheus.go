package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics may include it as a member to help them satisfy With semantics and save some code duplication.
type LabelValues []string

// With validates the input, and returns a new aggregate labelValues.
func (lvs LabelValues) With(labelValues ...string) LabelValues {
	if len(labelValues)%2 != 0 {
		labelValues = append(labelValues, "unknown")
	}
	return append(lvs, labelValues...)
}

// PrometheusCounter implements Counter, via a Prometheus CounterVec.
type PrometheusCounter struct {
	cv  *prometheus.CounterVec
	lvs LabelValues
}

// NewCounterFrom constructs and registers a Prometheus CounterVec,
// and returns a usable PrometheusCounter object.
func NewCounterFrom(opts prometheus.CounterOpts, labelNames []string) *PrometheusCounter {
	cv := prometheus.NewCounterVec(opts, labelNames)
	prometheus.MustRegister(cv)
	return NewCounter(cv)
}

// NewCounter wraps the CounterVec and returns a usable PrometheusCounter object.
func NewCounter(cv *prometheus.CounterVec) *PrometheusCounter {
	return &PrometheusCounter{
		cv: cv,
	}
}

// With implements PrometheusCounter.
func (c *PrometheusCounter) With(labelValues ...string) Counter {
	return &PrometheusCounter{
		cv:  c.cv,
		lvs: c.lvs.With(labelValues...),
	}
}

// Add implements Counter.
func (c *PrometheusCounter) Add(delta float64) {
	c.cv.With(makeLabels(c.lvs...)).Add(delta)
}

// PrometheusGauge implements Gauge, via a Prometheus GaugeVec.
type PrometheusGauge struct {
	gv  *prometheus.GaugeVec
	lvs LabelValues
}

// NewGaugeFrom constructs and registers a Prometheus GaugeVec,
// and returns a usable PrometheusGauge object.
func NewGaugeFrom(opts prometheus.GaugeOpts, labelNames []string) *PrometheusGauge {
	gv := prometheus.NewGaugeVec(opts, labelNames)
	prometheus.MustRegister(gv)
	return NewGauge(gv)
}

// NewGauge wraps the GaugeVec and returns a usable PrometheusGauge object.
func NewGauge(gv *prometheus.GaugeVec) *PrometheusGauge {
	return &PrometheusGauge{
		gv: gv,
	}
}

// With implements Gauge.
func (g *PrometheusGauge) With(labelValues ...string) Gauge {
	return &PrometheusGauge{
		gv:  g.gv,
		lvs: g.lvs.With(labelValues...),
	}
}

// Set implements Gauge.
func (g *PrometheusGauge) Set(value float64) {
	g.gv.With(makeLabels(g.lvs...)).Set(value)
}

// Add is supported by Prometheus GaugeVecs.
func (g *PrometheusGauge) Add(delta float64) {
	g.gv.With(makeLabels(g.lvs...)).Add(delta)
}

// Summary implements Histogram, via a Prometheus SummaryVec. The difference
// between a Summary and a Histogram is that Summaries don't require predefined
// quantile buckets, but cannot be statistically aggregated.
type Summary struct {
	sv  *prometheus.SummaryVec
	lvs LabelValues
}

// NewSummaryFrom constructs and registers a Prometheus SummaryVec,
// and returns a usable Summary object.
func NewSummaryFrom(opts prometheus.SummaryOpts, labelNames []string) *Summary {
	sv := prometheus.NewSummaryVec(opts, labelNames)
	prometheus.MustRegister(sv)
	return NewSummary(sv)
}

// NewSummary wraps the SummaryVec and returns a usable Summary object.
func NewSummary(sv *prometheus.SummaryVec) *Summary {
	return &Summary{
		sv: sv,
	}
}

// With implements Histogram.
func (s *Summary) With(labelValues ...string) Histogram {
	return &Summary{
		sv:  s.sv,
		lvs: s.lvs.With(labelValues...),
	}
}

// Observe implements Histogram.
func (s *Summary) Observe(value float64) {
	s.sv.With(makeLabels(s.lvs...)).Observe(value)
}

// PrometheusHistogram implements Histogram via a Prometheus HistogramVec. The difference
// between a Histogram and a Summary is that Histograms require predefined
// quantile buckets, and can be statistically aggregated.
type PrometheusHistogram struct {
	hv  *prometheus.HistogramVec
	lvs LabelValues
}

// NewHistogramFrom constructs and registers a Prometheus HistogramVec,
// and returns a usable PrometheusHistogram object.
func NewHistogramFrom(opts prometheus.HistogramOpts, labelNames []string) *PrometheusHistogram {
	hv := prometheus.NewHistogramVec(opts, labelNames)
	prometheus.MustRegister(hv)
	return NewHistogram(hv)
}

// NewHistogram wraps the HistogramVec and returns a usable PrometheusHistogram object.
func NewHistogram(hv *prometheus.HistogramVec) *PrometheusHistogram {
	return &PrometheusHistogram{
		hv: hv,
	}
}

// With implements Histogram.
func (h *PrometheusHistogram) With(labelValues ...string) Histogram {
	return &PrometheusHistogram{
		hv:  h.hv,
		lvs: h.lvs.With(labelValues...),
	}
}

// Observe implements Histogram.
func (h *PrometheusHistogram) Observe(value float64) {
	h.hv.With(makeLabels(h.lvs...)).Observe(value)
}

func makeLabels(labelValues ...string) prometheus.Labels {
	labels := prometheus.Labels{}
	for i := 0; i < len(labelValues); i += 2 {
		labels[labelValues[i]] = labelValues[i+1]
	}
	return labels
}
